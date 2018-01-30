package syncthing

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jpillora/overseer/fetcher"
	log "github.com/sirupsen/logrus"
)

// syncthing renames the OS for mac to macosx instead of darwin.
func matchRelease(filename string) bool {
	os := runtime.GOOS
	if os == "darwin" {
		os = "macosx"
	}

	return strings.Contains(filename, os) &&
		strings.Contains(filename, runtime.GOARCH)
}

func saveBinary(tarReader *tar.Reader, path string) error { //nolint interfacer
	dir := filepath.Dir(path)
	if _, statErr := os.Stat(dir); os.IsNotExist(statErr) {
		if mkdirErr := os.Mkdir(dir, 0700); mkdirErr != nil {
			return mkdirErr
		}
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0500)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, tarReader); err != nil {
		return err
	}

	log.Debug("wrote syncthing binary")

	return nil
}

// Fetch pulls down the latest syncthing binary to the provided path.
func Fetch(path string) error {
	f := &fetcher.Github{
		User:  "syncthing",
		Repo:  "syncthing",
		Asset: matchRelease,
	}

	if err := f.Init(); err != nil {
		return err
	}

	log.Debug("fetching new syncthing binary")

	gzReader, err := f.Fetch()
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()

		if err != nil {
			return err
		}

		// There are config files that are named the same thing as the binary. As
		// they're in etc directories, ignore those too.
		if strings.HasSuffix(header.Name, "/syncthing") &&
			!strings.Contains(header.Name, "/etc/") {
			return saveBinary(tarReader, path)
		}
	}
}
