package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// TODO: update the usage instructions
	statusHelp = `
    Show the status of existing sync tasks.
    `

	statusCmd = &cobra.Command{
		Use:   "status [flags]",
		Short: "Show the status of existing sync tasks.",
		Long:  statusHelp,
		Run:   runStatus,
		// TODO: BashCompletionFunction
	}
)

func runStatus(_ *cobra.Command, args []string) {
	log.Debug("status")
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
