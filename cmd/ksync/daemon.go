package main

import (
	"path/filepath"

	daemon "github.com/sevlyar/go-daemon"

	"github.com/ksync/ksync/pkg/cli"
)

func getDaemonContext() *daemon.Context {
	rootDir := cli.ConfigPath()
	return &daemon.Context{
		PidFileName: filepath.Join(rootDir, "daemon.pid"),
		LogFileName: filepath.Join(rootDir, "daemon.log"),
		WorkDir:     rootDir}
}
