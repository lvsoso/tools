package tasks

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/git-lfs/git-lfs/v3/lfs"
	"github.com/stretchr/testify/assert"
)

func Test_clone(t *testing.T) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "lfs-migrate-")
	assert.Nil(t, err)
	defer func() {
		os.RemoveAll(tmpDir)
	}()

	type args struct {
		workDir   string
		sourceUri string
		commit    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "plain repo",
			args: args{
				workDir:   filepath.Join(tmpDir, "clone"),
				sourceUri: "https://gitea.com/lvoooo/init.git",
				commit:    "70a233edcae446a8babbf29f21d05b829afc189d",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := clone(tt.args.workDir, tt.args.sourceUri, tt.args.commit); (err != nil) != tt.wantErr {
				t.Errorf("clone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	t.Log("finished")
}

func Test_filter(t *testing.T) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "lfs-migrate-")
	assert.Nil(t, err)
	defer func() {
		os.RemoveAll(tmpDir)
	}()

	target := filepath.Join(tmpDir, "clone")
	err = clone(
		target,
		"https://gitea.com/lvoooo/init.git",
		"70a233edcae446a8babbf29f21d05b829afc189d",
	)
	assert.Nil(t, err)

	type args struct {
		workDir string
	}
	tests := []struct {
		name    string
		args    args
		want    []MigrateObj
		wantErr bool
	}{
		{
			name: "plain repo",
			args: args{
				workDir: target,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getObjs(tt.args.workDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("filter() = %v, want %v", got, tt.want)
			}
		})
	}
	t.Log("finished")
}

func Test_decode(t *testing.T) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "lfs-migrate-")
	assert.Nil(t, err)
	defer func() {
		os.RemoveAll(tmpDir)
	}()

	target := filepath.Join(tmpDir, "clone")
	err = clone(
		target,
		"https://gitea.com/lvoooo/init.git",
		"45e14e89c73640810e688f3a3fdfb23142613792",
	)
	assert.Nil(t, err)

	objs, err := getObjs(target)
	assert.Nil(t, err)

	for _, o := range objs {
		t.Run(o.LocalPath, func(t *testing.T) {
			_, err := decode(o)
			if strings.HasPrefix(filepath.Base(o.LocalPath), "boot") {
				assert.Nil(t, err)
			} else {
				assert.ErrorIs(t, err, ErrNotALfsFile)
			}
		})
	}

	t.Log("")
}

func Test_migrate(t *testing.T) {
	m := MigrateObj{
		TargetPath: "/lfs/6e/12/6e12819814f936cdf125f2ed4640cf78a24a927d38916f257901d1d25e4a2a8d",
		p: &lfs.Pointer{
			Size: 524288000,
			Oid:  "6e12819814f936cdf125f2ed4640cf78a24a927d38916f257901d1d25e4a2a8d",
		},
		// TargetPath: "/lfs/c9/9b/c99b1345eb77db24c06b6e46763f8bc5199ba3e3e172c3f09c07b4024b8c14a6",
		// p: &lfs.Pointer{
		// 	Size: 10485760,
		// 	Oid:  "c99b1345eb77db24c06b6e46763f8bc5199ba3e3e172c3f09c07b4024b8c14a6",
		// },
	}
	gitUrl := "https://gitea.com/lvoooo/init.git"
	_, err := migrate(m, gitUrl)
	assert.Nil(t, err)
}
