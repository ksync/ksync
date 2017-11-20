package ksync

import (
	"github.com/cavaliercoder/grab"
	log "github.com/sirupsen/logrus"

	"github.com/vapor-ware/ksync/pkg/cli"
)

var (
	mirrorURL = "http://repo.joist.ws/mirror-all.jar"
)

// FetchMirror downloads the mirror jar into the correct location.
func FetchMirror() error {
	log.Debug("downloading mirror")

	client := grab.NewClient()
	req, err := grab.NewRequest(cli.ConfigPath(), mirrorURL)
	if err != nil {
		return err
	}

	resp := client.Do(req)
	<-resp.Done

	if err := resp.Err(); err != nil {
		return err
	}

	log.Debug("downloaded mirror")
	return nil
}
