// +build windows

package main

import (
	log "github.com/sirupsen/logrus"
)

func (c *cleanCmd) cleanLocal() {
	log.Info("Daemonization not supported on windows.")
}
