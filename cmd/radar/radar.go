package main

import (
	// homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	// "github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/cli"
)

var (
	globalUsage = `Map container names to local filesystem locations.`

	RootCmd = &cobra.Command{
		Use:              "radar",
		Short:            "Map container names to local filesystem locations.",
		Long:             globalUsage,
		PersistentPreRun: func(cmd *cobra.Command, args []string) { cli.InitLogging() },
	}
)

func main() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}

func init() {
	cobra.OnInitialize(func() { cli.InitConfig("radar") })

	cli.DefaultFlags(RootCmd, "radar")
}
