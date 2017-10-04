package input

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type Paths struct {
	Local  string
	Remote string
}

func GetPaths(args []string) Paths {
	return Paths{
		args[0],
		args[1],
	}
}

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
