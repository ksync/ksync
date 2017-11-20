package main

import (
	// "os"
	"runtime"
	// "time"
	// "text/template"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	// "github.com/inconshreveable/go-update"
	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"

	"github.com/vapor-ware/ksync/pkg/cli"
)

const (
	repoUsername = "vapor-ware"
	repoName     = "ksync"
)

type updateCmd struct {
	cli.BaseCmd
}

func (u *updateCmd) new() *cobra.Command {
	long := `Check for updates.`
	example := ``

	u.Init("ksync", &cobra.Command{
		Use:     "update",
		Short:   "Check for updates.",
		Long:    long,
		Example: example,
		Run:     u.run,
	})

	return u.Cmd
}

func validateOverseer() bool {
	log.Debugf("Checking if overseer is compatible with %s/%s", runtime.GOOS, runtime.GOARCH)
	if supported := overseer.IsSupported(); supported != true {
		log.Fatal("Overseer not compatible with this os or architecture")
		return supported
	}
	return true
}

func UpdateCheck() {
	overseer.Run(overseer.Config{
		Required: true,
		Program: runUpdater,
		Address: ":0000",
		NoRestart: true,
		Debug: true,
		Fetcher: &fetcher.Github{
			User: "vapor-ware",
			Repo: "ksync`",
		},
	})
}

func runUpdater(state overseer.State) {
	if ! validateOverseer() {
		log.Fatal("Update check failed")
	}
	log.Debug(state)
}

func (u *updateCmd) run(cmd *cobra.Command, args []string) {
	UpdateCheck()
}
