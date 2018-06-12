// +build windows

package ksync

import (
	"path/filepath"

	"github.com/vapor-ware/ksync/pkg/cli"
)

func (s *Syncthing) binPath() string {
	return filepath.Join(cli.ConfigPath(), "bin", "syncthing.exe")
}
