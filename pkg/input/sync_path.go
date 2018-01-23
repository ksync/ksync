package input

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/phayes/permbits"
)

// SyncPath has both the local and remote file paths for a specific sync.
type SyncPath struct {
	Local  string
	Remote string
}

// GetSyncPath constructs a SyncPath from command line arguments.
func GetSyncPath(args []string) SyncPath {
	return SyncPath{
		args[0],
		args[1],
	}
}

// localPathHasPermission checks a given root directory, and all children, for
// `rw` permissions for the current user.
func (s *SyncPath) localPathHasPermission() error {
	root, err := filepath.Abs(s.Local)
	if err != nil {
		return err
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		permissions, errStat := permbits.Stat(path)
		if errStat != nil {
			return errStat
		}

		switch {
		case !permissions.UserRead():
			return fmt.Errorf("File %s is not readable. It is set to %v", path, permissions)
		case !permissions.UserWrite():
			return fmt.Errorf("File %s is not writable. It is set to %v", path, permissions)
		}

		return nil
	})
	return err
}

// Validator ensures the SyncPath is valid and can be used to configure a
// sync.
func (s *SyncPath) Validator() error {
	if s.Local == "" {
		return fmt.Errorf("must specify a local path")
	}

	if s.Remote == "" {
		return fmt.Errorf("must specify a remote path")
	}

	if !filepath.IsAbs(s.Local) {
		return fmt.Errorf("local path must be absolute")
	}

	if !filepath.IsAbs(s.Remote) {
		return fmt.Errorf("remote path must be absolute")
	}

	if err := s.localPathHasPermission(); err != nil {
		return err
	}

	return nil
}
