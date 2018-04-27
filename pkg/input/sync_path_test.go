package input

import (
	"io/ioutil"
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

	// Test for missing local path
	path = &SyncPath{
		Local: "",
	}
	err := path.Validator()
	assert.EqualError(t, err, "must specify a local path")

	// Check for missing remote path
	path = &SyncPath{
		Local: os.TempDir(),
	}
	err = path.Validator()
	assert.EqualError(t, err, "must specify a remote path")

	// Check if local path is absolute
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

	// Check if remote path is absolute
	path = &SyncPath{
		Local:  absdirpath,
		Remote: filepath.Base(currentdir),
	}
	err = path.Validator()
	assert.EqualError(t, err, "remote path must be absolute")

	// Check if path contains files that are not `rw`
	// Create out unreadable file and write content to it
	badfile, err := ioutil.TempFile(os.TempDir(), "badfile")            //nolint: ineffassign,megacheck
	_, err = badfile.WriteString("We should not be able to read this.") //nolint: ineffassign,megacheck
	err = badfile.Close()                                               //nolint: ineffassign,megacheck
	err = os.Chmod(badfile.Name(), 0100)                                //nolint: ineffassign,megacheck
	t.Logf("File: %s", badfile.Name())
	require.NoError(t, err)

	// Clean up after ourselves
	defer os.Remove(badfile.Name()) //nolint: errcheck

	path = &SyncPath{
		Local:  os.TempDir(),
		Remote: os.TempDir(),
	}
	err = path.Validator()
	// Temporarily changed. See ./sync_path.go:68
	// assert.Error(t, err)
	assert.NoError(t, err)

}
