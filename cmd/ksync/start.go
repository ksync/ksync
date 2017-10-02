package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// TODO: update the usage instructions
	startHelp = `
    Start syncing between a local and remote directory.
    `

	startCmd = &cobra.Command{
		Use:   "start [flags] [local path] [remote path]",
		Short: "Start syncing between a local and remote directory.",
		Long:  listHelp,
		Run:   runStart,
		// TODO: BashCompletionFunction
	}
)

func runStart(_ *cobra.Command, args []string) {
	log.Debug("start")
}

func init() {
	RootCmd.AddCommand(startCmd)
}
