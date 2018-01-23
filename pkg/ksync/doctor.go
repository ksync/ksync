package ksync

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/vapor-ware/ksync/pkg/cli"
)

// JavaBin contains the name of the java binary
var JavaBin = "java"

// MirrorBin contains the name of the mirror binary
var MirrorBin = "mirror-all.jar"

// HasJava checks if java is on the path.
func HasJava() bool {
	if _, err := exec.LookPath(JavaBin); err != nil {
		return false
	}

	return true
}

// HasMirror verifies that mirror has been downloaded and placed in the right
// location
func HasMirror() bool {
	if _, err := os.Stat(
		filepath.Join(cli.ConfigPath(), MirrorBin)); err != nil {
		return false
	}

	return true
}
