package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
	"github.com/vapor-ware/ksync/pkg/ksync/cluster"
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
		(&cleanCmd{}).new(),
		(&createCmd{}).new(),
		(&deleteCmd{}).new(),
		(&getCmd{}).new(),
		(&initCmd{}).new(),
		(&watchCmd{}).new(),
		(&versionCmd{}).new(),
		(&updateCmd{}).new(),
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

	flags.String(
		"image",
		fmt.Sprintf("vaporio/ksync:git-%s", ksync.GitCommit),
		"the image to use for radar.")
	if err := flags.MarkHidden("image"); err != nil {
		log.Fatal(err)
	}

	if err := cli.BindFlag(
		viper.GetViper(), flags.Lookup("image"), "ksync"); err != nil {

		log.Fatal(err)
	}

	flags.Int(
		"port",
		40322,
		"port on which the server will listen")

	if err := cli.BindFlag(
		viper.GetViper(), flags.Lookup("port"), "ksync"); err != nil {

		log.Fatal(err)
	}

	flags.String(
		"apikey",
		"ksync",
		"api key used for authentication with syncthing")

	if err := cli.BindFlag(
		viper.GetViper(), flags.Lookup("apikey"), "ksync"); err != nil {

		log.Fatal(err)
	}

	flags.Int(
		"syncthing-port",
		8384,
		"port on which the syncthing server will listen")

	if err := cli.BindFlag(
		viper.GetViper(), flags.Lookup("syncthing-port"), "ksync"); err != nil {

		log.Fatal(err)
	}
}

// TODO: dependencies should verify that they're usable
// (and return errors otherwise).
func initPersistent(cmd *cobra.Command, args []string) {
	cli.InitLogging()

	initKubeClient()

	cluster.SetImage(viper.GetString("image"))
}

func initKubeClient() {
	err := cluster.InitKubeClient(viper.GetString("context"))
	if err != nil {
		log.Fatalf("Error creating kubernetes client: %v", err)
	}
}
