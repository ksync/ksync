package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// TODO: update the usage instructions
	stopHelp = `
    Stop syncing between a local and remote directory.
    `

	stopCmd = &cobra.Command{
		Use:   "stop [flags] [id]",
		Short: "Stop syncing between a local and remote directory.",
		Long:  stopHelp,
		Run:   runStop,
		// TODO: BashCompletionFunction
	}
)

func runStop(_ *cobra.Command, args []string) {
	log.Debug("stop")
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
