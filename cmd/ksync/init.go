package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/ksync"
)

var (
	// TODO: update the usage instructions
	initHelp = `
    Prepare the cluster.
    `

	initCmd = &cobra.Command{
		Use:   "init [flags]",
		Short: "Prepare the cluster.",
		Long:  initHelp,
		Run:   runInit,
	}
)

// TODO: upgrade currently doesn't work because the template doesn't change
// (when on canary).
func runInit(_ *cobra.Command, _ []string) {
	ksync.InitRadar(viper.GetBool("upgrade"))

	log.Print("wee")
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP(
		"upgrade",
		"u",
		false,
		"Upgrade the currently running version.")

	viper.BindPFlag("upgrade", initCmd.Flags().Lookup("upgrade"))
}
