package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
)

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

func main() {
	RootCmd.AddCommand(
		(&createCmd{}).new(),
		(&deleteCmd{}).new(),
		(&getCmd{}).new(),
		(&initCmd{}).new(),
		(&listCmd{}).new(),
		(&runCmd{}).new(),
		(&watchCmd{}).new(),
	)
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}

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
