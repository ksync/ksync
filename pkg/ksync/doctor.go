package ksync

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/vapor-ware/ksync/pkg/cli"
)

// HasJava checks if java is on the path.
func HasJava() bool {
	if _, err := exec.LookPath("java"); err != nil {
		return false
	}

	return true
}

// HasMirror verifies that mirror has been downloaded and placed in the right
// location
func HasMirror() bool {
	if _, err := os.Stat(
		filepath.Join(cli.ConfigPath(), "mirror-all.jar")); err != nil {
		return false
	}

	return true
}
