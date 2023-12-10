package lfscache

import (
	"crypto/sha256"
	"demo-tui/internal/gitop"
	"demo-tui/log"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var logger = log.Logger

const (
	cacheObjects = ".git/lfs/objects"
	cacheTmp     = ".git/lfs/tmp"
)

const lfsPointerTemplate = `version https://git-lfs.github.com/spec/v1
oid sha256:%s
size %d
`

type Object struct {
	name      string
	src       string
	dst       string
	oid       string
	size      int64
	cachePath string
	tmp       string
}

func initDir(repoRootDir string) error {
	// create .git/lfs/objects
	objectsDir := filepath.Join(repoRootDir, cacheObjects)
	if _, err := os.Stat(objectsDir); os.IsNotExist(err) {
		err := os.MkdirAll(objectsDir, 0755)
		if err != nil {
			return err
		}
	}

	// create .git/lfs/tmp
	tmpDir := filepath.Join(repoRootDir, cacheTmp)
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		err := os.MkdirAll(tmpDir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func process(srcfile string, repoRootDir string, obj *Object, wg *sync.WaitGroup, concurrent chan int) {
	defer wg.Done()
	defer func() { <-concurrent }()

	// tmp file
	tmp, err := ioutil.TempFile(filepath.Join(repoRootDir, cacheTmp), "")
	if err != nil {
		panic(err)
	}
	defer tmp.Close()
	logger.Info(tmp.Name())

	// copy
	oidHash := sha256.New()
	writer := io.MultiWriter(oidHash, tmp)
	reader, err := os.Open(srcfile)
	if err != nil {
		panic(err)
	}

	size, err := io.Copy(writer, reader)
	if err != nil {
		panic(err)
	}
	logger.Info(size)

	oid := hex.EncodeToString(oidHash.Sum(nil))
	logger.Info(oid)

	obj.name = filepath.Base(srcfile)
	obj.src = srcfile
	obj.oid = oid
	obj.size = size
	obj.tmp = tmp.Name()
	obj.cachePath = filepath.Join(repoRootDir, cacheObjects, oid[:2], oid[2:4], oid)
}

func LfsCache(sourceDir string, repoRootDir string, parallels int) error {
	if sourceDir == "" || repoRootDir == "" {
		return errors.New("empty source dir or repo root url")
	}

	// find .gitattributes
	attrFile := filepath.Join(repoRootDir, ".gitattributes")
	_, err := os.Stat(attrFile)
	if os.IsNotExist(err) {
		return errors.New("not found .gitattributes")
	} else if err != nil {
		return err
	}

	matcher, err := gitop.NewMatcher(attrFile)
	if err != nil {
		return err
	}

	// lfs  file filter rule ?
	err = initDir(repoRootDir)
	if err != nil {
		return err
	}

	srcFiles := []string{}
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			srcFiles = append(srcFiles, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	objests := make([]Object, len(srcFiles))
	concurrent := make(chan int, parallels)

	for idx, srcfile := range srcFiles {
		dstPath := filepath.Join(repoRootDir, strings.Replace(srcfile, sourceDir, "", 1))
		logger.Infof("dstPath: %s", dstPath)
		matched, err := matcher.MatchLfs(dstPath)
		if err != nil {
			return err
		}
		if !matched {
			logger.Infof("dstPath: %s            !matched", dstPath)
			continue
		}

		concurrent <- idx
		wg.Add(1)
		objests[idx].dst = dstPath
		go process(srcfile, repoRootDir, &objests[idx], &wg, concurrent)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	for _, o := range objests {
		logger.Infof("%+v", o)
		if o.dst == "" {
			logger.Infof("srcPath: %s   skipped", o.src)
			continue
		}

		if _, err := os.Stat(o.dst); os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(o.dst), 0755)
			if err != nil {
				return err
			}
		}

		// write pointer
		pointer := fmt.Sprintf(lfsPointerTemplate, o.oid, o.size)
		f, err := os.OpenFile(o.dst, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString(pointer)
		if err != nil {
			return err
		}

		// rename file
		if _, err := os.Stat(o.cachePath); os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(o.cachePath), 0755)
			if err != nil {
				return err
			}
		}
		err = os.Rename(o.tmp, o.cachePath)
		if err != nil {
			return err
		}

	}

	// repoFiles := []string{}
	return nil
}
