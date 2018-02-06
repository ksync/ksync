// +build !windows

package main

import (
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
	daemon "github.com/timfallmk/go-daemon"
)

func (c *cleanCmd) cleanLocal() {
	context := getDaemonContext()
	if _, err := context.Search(); err != nil {
		log.Infoln("No daemonized process found. Nothing to clean locally.")
		log.Warningln(err)
		return
	}

	pid, err := daemon.ReadPidFile(context.PidFileName)
	if err != nil {
		log.Fatal(err)
	}

	if err := syscall.Kill(-pid, os.Interrupt.(syscall.Signal)); err != nil {
		log.Fatal(err)
	}

	// Clean up after the process since it seems incapable of doing that itself
	if err := os.Remove(context.PidFileName); err != nil {
		log.Fatal(err)
	}
}
