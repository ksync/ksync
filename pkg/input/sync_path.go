package input

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// SyncPath has both the local and remote file paths for a specific sync.
type SyncPath struct {
	Local  string
	Remote string
}

// GetPaths constructs a SyncPath from command line arguments.
func GetSyncPath(args []string) SyncPath {
	return SyncPath{
		args[0],
		args[1],
	}
}

// Validator ensures the SyncPath is valid and can be used to configure a
// sync.
func (s *SyncPath) Validator() {
	if s.Local == "" {
		log.Fatal("Must specify a local path")
	}

	if s.Remote == "" {
		log.Fatal("Must specify a remote path")
	}

	if !filepath.IsAbs(s.Local) {
		log.Fatal("Local path must be absolute.")
	}

	if _, err := os.Stat(s.Local); os.IsNotExist(err) {
		log.Fatal("Local path must exist.")
	}

	if !filepath.IsAbs(s.Remote) {
		log.Fatal("Remote path must be absolute.")
	}
}
