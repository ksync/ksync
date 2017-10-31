package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/cli"
)

var (
	globalUsage = `Map container names to local filesystem locations.`

	rootCmd = &cobra.Command{
		Use:   "radar",
		Short: "Map container names to local filesystem locations.",
		Long:  globalUsage,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cli.InitLogging()
		},
	}
)

// Main runs the server instance
func main() {
	rootCmd.AddCommand(
		(&serveCmd{}).new(),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}

// Init initializes the server instance
func init() {
	cobra.OnInitialize(func() {
		if err := cli.InitConfig("radar"); err != nil {
			log.Fatal(err)
		}
	})

	if err := cli.DefaultFlags(rootCmd, "radar"); err != nil {
		log.Fatal(err)
	}

	flags := rootCmd.PersistentFlags()
	flags.String(
		"pod-name",
		"",
		"Name of the pod this is running inside.")
	if err := cli.BindFlag(
		viper.GetViper(), flags.Lookup("pod-name"), "radar"); err != nil {

		log.Fatal(err)
	}
}
