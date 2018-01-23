package main

import (
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/jpillora/overseer/fetcher"
	"github.com/timfallmk/overseer"

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
		Hidden:  false,
	})

	return u.Cmd
}

func validateOverseer() bool {
	log.Debugf("Checking if overseer is compatible with %s/%s", runtime.GOOS, runtime.GOARCH)
	if !overseer.IsSupported() {
		log.Fatal("Overseer not compatible with this os or architecture")
		return overseer.IsSupported()
	}
	return true
}

// UpdateCheck is the wrapping function that launches the overseer process and
// monitors the child process. In this case it just runs the update check and
// quits.
func UpdateCheck() {
	overseer.Run(overseer.Config{
		Required:  true,
		Program:   runUpdater,
		Address:   ":0000",
		NoRestart: true,
		Debug:     true,
		Fetcher: &fetcher.Github{
			User: repoUsername,
			Repo: repoName,
		},
	})
}

func runUpdater(state overseer.State) {
	if !validateOverseer() {
		log.Fatal("Update check failed")
	}
	log.Debug(state)
}

func (u *updateCmd) run(cmd *cobra.Command, args []string) {
	UpdateCheck()
}
