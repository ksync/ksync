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
		(&watchCmd{}).new(),
		(&versionCmd{}).new(),
	)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}

func init() {
	cobra.OnInitialize(func() {
		if err := cli.InitConfig("ksync"); err != nil {
			log.Fatal(err)
		}
	})

	if err := cli.DefaultFlags(rootCmd, "ksync"); err != nil {
		log.Fatal(err)
	}

	flags := rootCmd.PersistentFlags()
	flags.StringP(
		"namespace",
		"n",
		"default",
		"namespace to use.")
	if err := cli.BindFlag(
		viper.GetViper(), flags.Lookup("namespace"), "ksync"); err != nil {

		log.Fatal(err)
	}

	flags.String(
		"context",
		"",
		"name of the kubeconfig context to use")
	if err := cli.BindFlag(
		viper.GetViper(), flags.Lookup("context"), "ksync"); err != nil {

		log.Fatal(err)
	}

	// TODO: can this be hidden?
	flags.String(
		"image",
		"gcr.io/elated-embassy-152022/ksync/ksync:canary",
		// TODO: this help text could be way better
		"the image to use for running things locally.")
	if err := cli.BindFlag(
		viper.GetViper(), flags.Lookup("image"), "ksync"); err != nil {

		log.Fatal(err)
	}
}

// TODO: dependencies should verify that they're usable
// (and return errors otherwise).
func initPersistent(cmd *cobra.Command, args []string) {
	cli.InitLogging()

	initKubeClient()

	ksync.SetImage(viper.GetString("image"))
}

func initKubeClient() {
	err := ksync.InitKubeClient(viper.GetString("context"))
	if err != nil {
		log.Fatalf("Error creating kubernetes client: %v", err)
	}
}
