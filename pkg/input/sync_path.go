package input

import (
	"fmt"
	"path/filepath"
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

	return nil
}
