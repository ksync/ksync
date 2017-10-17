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

	rootCmd = &cobra.Command{
		Use:              "ksync",
		Short:            "Inspect and sync files from remote containers.",
		Long:             globalUsage,
		PersistentPreRun: initPersistent,
	}
)

func main() {
	rootCmd.AddCommand(
		(&createCmd{}).new(),
		(&deleteCmd{}).new(),
		(&getCmd{}).new(),
		(&initCmd{}).new(),
		(&listCmd{}).new(),
		(&runCmd{}).new(),
		(&watchCmd{}).new(),
	)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}

func init() {
	cobra.OnInitialize(func() { cli.InitConfig("ksync") })

	cli.DefaultFlags(rootCmd, "ksync")

	rootCmd.PersistentFlags().StringP(
		"namespace",
		"n",
		"default",
		"namespace to use.")

	viper.BindPFlag("namespace", rootCmd.PersistentFlags().Lookup("namespace"))
	viper.BindEnv("namespace", "KSYNC_NAMESPACE")

	rootCmd.PersistentFlags().String(
		"context",
		"",
		"name of the kubeconfig context to use")

	viper.BindPFlag("context", rootCmd.PersistentFlags().Lookup("context"))
	viper.BindEnv("context", "KSYNC_CONTEXT")
}

func initPersistent(cmd *cobra.Command, args []string) {
	cli.InitLogging()
	initKubeClient()
	initDockerClient()
}

func initKubeClient() {
	err := ksync.InitKubeClient(viper.GetString("context"), viper.GetString("namespace"))
	if err != nil {
		log.Fatalf("Error creating kubernetes client: %v", err)
	}
}

// TODO: should this be scoped only to commands that use docker?
func initDockerClient() {
	err := ksync.InitDockerClient()
	if err != nil {
		log.Fatalf("Error creating docker client: %v", err)
	}
}
