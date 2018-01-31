package main

import (
	"path/filepath"

	daemon "github.com/timfallmk/go-daemon"

	"github.com/vapor-ware/ksync/pkg/cli"
)

func getDaemonContext() *daemon.Context {
	rootDir := cli.ConfigPath()
	return &daemon.Context{
		PidFileName: filepath.Join(rootDir, "daemon.pid"),
		LogFileName: filepath.Join(rootDir, "daemon.log"),
		WorkDir:     rootDir}
}
