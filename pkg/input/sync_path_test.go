package input

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSyncPath(t *testing.T) {
	args := []string{"argle", "bargle"}
	path := GetSyncPath(args)

	assert.NotEmpty(t, path)
	assert.Equal(t, path, SyncPath{Local: "argle", Remote: "bargle"})
	assert.NotEqual(t, path.Local, "fargle")
}

func TestValidator(t *testing.T) {
	var path = new(SyncPath)
	path = &SyncPath{
		Local: "",
	}
	err := path.Validator()
	assert.EqualError(t, err, "must specify a local path")

	path = &SyncPath{
		Local: os.TempDir(),
	}
	err = path.Validator()
	assert.EqualError(t, err, "must specify a remote path")

	currentdir, patherr := os.Getwd()
	require.NoError(t, patherr)
	absdirpath, abspatherr := filepath.Abs(currentdir)
	require.NoError(t, abspatherr)

	path = &SyncPath{
		Local:  filepath.Base(currentdir),
		Remote: os.TempDir(),
	}
	err = path.Validator()
	assert.EqualError(t, err, "local path must be absolute")

	path = &SyncPath{
		Local:  absdirpath,
		Remote: filepath.Base(currentdir),
	}
	err = path.Validator()
	assert.EqualError(t, err, "remote path must be absolute")

}
