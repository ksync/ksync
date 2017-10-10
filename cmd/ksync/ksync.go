package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
)

// TODO: should there be an init command that lets you do a DaemonSet (instead of single pods)
var (
	// TODO: update usage with flags
	globalUsage = `Inspect and sync files from remote containers.`

	RootCmd = &cobra.Command{
		Use:              "ksync",
		Short:            "Inspect and sync files from remote containers.",
		Long:             globalUsage,
		PersistentPreRun: initPersistent,
	}
)

// Main creates the command line client and adds all subcommands
func main() {
	RootCmd.AddCommand(
		(&CreateCmd{}).New(),
		(&DeleteCmd{}).New(),
		(&GetCmd{}).New(),
		(&InitCmd{}).New(),
		(&ListCmd{}).New(),
		(&RunCmd{}).New(),
		(&WatchCmd{}).New(),
	)
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}

// Init initializes a new ksync client and populates it with the default
// commands and parameters
func init() {
	cobra.OnInitialize(func() { cli.InitConfig("ksync") })

	cli.DefaultFlags(RootCmd, "ksync")

	RootCmd.PersistentFlags().StringP(
		"namespace",
		"n",
		"default",
		"namespace to use.")

	viper.BindPFlag("namespace", RootCmd.PersistentFlags().Lookup("namespace"))
	viper.BindEnv("namespace", "KSYNC_NAMESPACE")

	RootCmd.PersistentFlags().String(
		"context",
		"",
		"name of the kubeconfig context to use")

	viper.BindPFlag("context", RootCmd.PersistentFlags().Lookup("context"))
	viper.BindEnv("context", "KSYNC_CONTEXT")
}

// InitPersistent initializes the functions that persitent between runs. This
// includes the server (radar) and logging options.
func initPersistent(cmd *cobra.Command, args []string) {
	cli.InitLogging()
	initClient()
	ksync.InitRadarOpts()
}

// InitClient initializes the kubernetes client options.
func initClient() {
	err := ksync.InitClient(viper.GetString("context"), viper.GetString("namespace"))
	if err != nil {
		log.Fatalf("Error creating kubernetes client: %v", err)
	}
}
