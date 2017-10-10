package input

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Paths defines the fields for local and remote paths
type Paths struct {
	Local  string
	Remote string
}

// GetPaths returns an array of paths from a Paths object
func GetPaths(args []string) Paths {
	return Paths{
		args[0],
		args[1],
	}
}

// Validator validates paths for correct syntax and existence
func (this *Paths) Validator() {
	if this.Local == "" {
		log.Fatal("Must specify a local path")
	}

	if this.Remote == "" {
		log.Fatal("Must specify a remote path")
	}

	if !filepath.IsAbs(this.Local) {
		log.Fatal("Local path must be absolute.")
	}

	if _, err := os.Stat(this.Local); os.IsNotExist(err) {
		log.Fatal("Local path must exist.")
	}

	if !filepath.IsAbs(this.Remote) {
		log.Fatal("Remote path must be absolute.")
	}
}
