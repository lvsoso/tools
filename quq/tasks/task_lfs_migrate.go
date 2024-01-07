package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	qconfig "quq/config"
	"quq/storage"

	lfserrors "github.com/git-lfs/git-lfs/v3/errors"
	"github.com/git-lfs/git-lfs/v3/lfs"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	lfshttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/hibiken/asynq"
)

var (
	ErrDecodePointerError     = errors.New("decode pointer error")
	ErrNotALfsFile            = errors.New("not a lfs file")
	ErrGetLfsFile             = errors.New("get lfs file failed")
	ErrGetLfsFileRead         = errors.New("get lfs file read error")
	ErrGetLfsFileMisMatchSize = errors.New("get lfs file size error")
	ErrGetStorageConfig       = errors.New("get storage config error")
	ErrCreateStorageFailed    = errors.New("create storage error")
)

const (
	TypeLfsMigrate = "tasks:lfs:migrate"
)

type LfsMigratePayload struct {
	SourceUri string
	Commit    string
	Target    string
}

func NewLfsMigratePayload(sourceUri string, commit string, target string) (*asynq.Task, error) {
	payload, err := json.Marshal(LfsMigratePayload{
		SourceUri: sourceUri,
		Commit:    commit,
		Target:    target,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeLfsMigrate, payload), nil
}

func HandleLfsMigratePayload(ctx context.Context, t *asynq.Task) error {
	processStep := make(chan int64, 100)
	errorQueue := make(chan interface{})
	var p LfsMigratePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	go handleLfsMigratePayload(ctx, t, p, errorQueue, processStep)
	var err interface{}
	for {
		select {
		case <-ctx.Done():
			Logger.Info("canceled")
			return nil
		case err = <-errorQueue:
			if err != nil {
				Logger.Error(err)
				return fmt.Errorf("unknown %v", err)
			}
			Logger.Info("finnished")
			return nil
		case step := <-processStep:
			Logger.Info(fmt.Sprintf("step %d \n", step))
		}
	}
}

func handleLfsMigratePayload(
	ctx context.Context,
	t *asynq.Task,
	p LfsMigratePayload,
	errorQueue chan interface{},
	processStep chan int64) {
	defer func() {
		if err := recover(); err != nil {
			Logger.Error(err)
			errorQueue <- err
		}
	}()

	count := 100
	for i := 0; i < int(count); i++ {
		time.Sleep(1 * time.Second)
		Logger.Info(fmt.Sprintf("Count=%d", i))
		select {
		case <-ctx.Done():
			Logger.Info("canceled recived")
			return
		default:
		}
		processStep <- int64(i)
	}

	tmpDir, err := os.MkdirTemp(os.TempDir(), "lfs-migrate")
	if err != nil {
		errorQueue <- err
	}
	defer func() {
		err = os.RemoveAll(tmpDir)
		if err != nil {
			Logger.Error(err.Error())
			return
		}
	}()

	// clone
	err = clone(tmpDir, p.SourceUri, p.Commit)
	if err != nil {
		Logger.Error(err.Error())
		return
	}

	// getObjs
	objs, err := getObjs(tmpDir)
	if err != nil {
		Logger.Error(err.Error())
		errorQueue <- err
		return
	}

	for _, o := range objs {
		Logger.Info("hanlde ...", o.LocalPath)
		// decode
		lfsO, err := decode(o)
		if err != nil {
			if !errors.Is(err, ErrNotALfsFile) {
				Logger.Error(err)
				errorQueue <- err
				return
			}
			continue
		}

		Logger.Info("Oid ...", lfsO.p.Oid, lfsO.p.Size)

		// // migrate
		// err = migrate(o)
		// if err != nil {
		// 	Logger.Error(err.Error())
		// 	errorQueue <- err
		// 	return
		// }
	}
	Logger.Info("finished")
}

func clone(workDir string, sourceUri string, commit string) error {
	repo, err := git.PlainInit(workDir, false)
	if err != nil {
		return err
	}

	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{sourceUri},
	})
	if err != nil {
		return err
	}
	var target = plumbing.NewBranchReferenceName("target")
	err = remote.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			config.RefSpec(plumbing.ReferenceName(commit) + ":" + target),
		},
		Depth: 1,
		Auth: &lfshttp.BasicAuth{
			Username: qconfig.Config.GitUsername,
			Password: qconfig.Config.GitPassword,
		},
	})
	if err != nil {
		return err
	}

	tree, err := repo.Worktree()
	if err != nil {
		return err
	}
	return tree.Checkout(&git.CheckoutOptions{
		Branch: target})
}

type MigrateObj struct {
	LocalPath  string
	TargetPath string
	Size       int64
	p          *lfs.Pointer
}

func getObjs(workDir string) ([]MigrateObj, error) {
	mo := []MigrateObj{}
	err := filepath.Walk(workDir, func(path string, info fs.FileInfo, err error) error {
		if strings.HasPrefix(path, filepath.Join(workDir, ".git")+"/") {
			Logger.Info("skip '.git' ...")
			return nil
		}
		if info.IsDir() {
			return nil
		}

		Logger.Info(path)

		mo = append(mo, MigrateObj{
			LocalPath: path,
			Size:      info.Size(),
		})

		return nil
	})
	if err != nil {
		return nil, err
	}
	return mo, nil
}

func decode(m MigrateObj) (MigrateObj, error) {
	if m.Size > 1024 {
		return m, ErrNotALfsFile
	}

	p, err := lfs.DecodePointerFromFile(m.LocalPath)
	if p != nil {
		m.p = p
		return m, nil
	}

	if lfserrors.IsNotAPointerError(err) {
		return m, ErrNotALfsFile
	}

	return m, err
}

func migrate(m MigrateObj, gitUrl string) (MigrateObj, error) {
	batchResp, err := GetLinkObj(m, gitUrl)
	if err != nil {
		return m, err
	}

	o := batchResp.Objects[0]
	req, err := http.NewRequest("GET", o.Actions.Download.Href, nil)
	if err != nil {
		Logger.Error(err)
		return m, err
	}
	for k, v := range o.Actions.Download.Header {
		req.Header.Set(k, v)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		Logger.Error(err)
		return m, err
	}
	defer func() {
		res.Body.Close()
	}()

	name := "azureblob"
	tc, err := qconfig.Config.GetConfig(name)
	if err != nil {
		Logger.Error(err)
		return m, ErrGetStorageConfig
	}

	st, err := storage.NewStorage(name, tc)
	if err != nil {
		Logger.Error(err)
		return m, ErrCreateStorageFailed
	}

	// f, err := os.OpenFile("/tmp/init/boot4.img", os.O_RDONLY, os.ModePerm)
	// if err != nil {
	// 	Logger.Error(err)
	// 	return m, ErrGetLfsFileRead
	// }
	// defer f.Close()

	// file, err := st.Put(m.TargetPath, f)
	file, err := st.Put(m.TargetPath, res.Body)
	// by, err := io.ReadAll(res.Body)
	if err != nil {
		Logger.Error(err)
		return m, ErrGetLfsFileRead
	}

	// if len(by) != int(m.p.Size) {
	// 	return m, ErrGetLfsFileMisMatchSize
	// }

	if int(file.Stat().Size()) != int(m.p.Size) {
		return m, ErrGetLfsFileMisMatchSize
	}

	return m, nil
}

func GetLinkObj(m MigrateObj, gitUrl string) (*BatchResp, error) {
	bg := BatchReq{
		Operation: "download",
		Transfers: []string{"basic"},
		Objects: []Object{{
			Oid:  m.p.Oid,
			Size: m.p.Size,
		}},
		HashAlgo: "sha256",
	}

	data, err := json.Marshal(bg)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	buf := bytes.NewReader(data)
	lfsUrl := gitUrl + "/info/lfs/objects/batch"
	req, err := http.NewRequest("POST", lfsUrl, buf)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.git-lfs+json")
	req.Header.Set("Content-Type", "application/vnd.git-lfs+json")

	req.SetBasicAuth(
		qconfig.Config.GitUsername,
		qconfig.Config.GitPassword)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	defer func() {
		res.Body.Close()
	}()

	by, err := io.ReadAll(res.Body)
	if err != nil {
		Logger.Error(err)
		return nil, ErrGetLfsFileRead
	}

	br := BatchResp{}
	err = json.Unmarshal(by, &br)
	if err != nil {
		Logger.Error(ErrGetLfsFileMisMatchSize)
		return nil, ErrGetLfsFileMisMatchSize
	}
	return &br, nil
}
