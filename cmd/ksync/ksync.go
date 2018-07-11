package main

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
	"github.com/vapor-ware/ksync/pkg/ksync/cluster"
	"github.com/vapor-ware/ksync/pkg/ksync/doctor"
)

var (
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
		(&doctorCmd{}).new(),
		(&getCmd{}).new(),
		(&initCmd{}).new(),
		(&reloadCmd{}).new(),
		(&watchCmd{}).new(),
		(&versionCmd{}).new(),
		(&updateCmd{}).new(),
	)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}

func localFlags(flags *pflag.FlagSet) {
	flags.StringP(
		"namespace",
		"n",
		"default",
		"namespace to use")
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

	flags.Int(
		"port",
		40322,
		"port on watch listens on locally")

	if err := cli.BindFlag(
		viper.GetViper(), flags.Lookup("port"), "ksync"); err != nil {

		log.Fatal(err)
	}

	flags.StringP(
		"output",
		"o",
		"pretty",
		"output format to use (e.g. \"json\")")

	if err := cli.BindFlag(viper.GetViper(), flags.Lookup("output"), "ksync"); err != nil {
		log.Fatal(err)
	}
}

func remoteFlags(flags *pflag.FlagSet) {
	flags.String(
		"image",
		fmt.Sprintf("vaporio/ksync:git-%s", ksync.GitCommit),
		"the image to use on the cluster")
	if err := flags.MarkHidden("image"); err != nil {
		log.Fatal(err)
	}

	if err := cli.BindFlag(
		viper.GetViper(), flags.Lookup("image"), "ksync"); err != nil {

		log.Fatal(err)
	}

	flags.String(
		"apikey",
		"ksync",
		"api key used for authentication with syncthing")
	if err := flags.MarkHidden("apikey"); err != nil {
		log.Fatal(err)
	}

	if err := cli.BindFlag(
		viper.GetViper(), flags.Lookup("apikey"), "ksync"); err != nil {

		log.Fatal(err)
	}

	flags.Int(
		"syncthing-port",
		8384,
		"port on which the syncthing server will listen")
	if err := flags.MarkHidden("syncthing-port"); err != nil {
		log.Fatal(err)
	}

	if err := cli.BindFlag(
		viper.GetViper(), flags.Lookup("syncthing-port"), "ksync"); err != nil {

		log.Fatal(err)
	}

	flags.String(
		"graph-root",
		"/var/lib/docker",
		"root directory of the docker storage (graph) driver")
	if err := flags.MarkHidden("graph-root"); err != nil {
		log.Fatal(err)
	}

	if err := cli.BindFlag(
		viper.GetViper(), flags.Lookup("graph-root"), "ksync"); err != nil {

		log.Fatal(err)
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
	localFlags(flags)
	remoteFlags(flags)
}

func initPersistent(cmd *cobra.Command, args []string) {
	cli.InitLogging()

	// This is a super special case where we don't want to initialize the k8s
	// client, instead waiting to test it as part of the doctor process.
	if !strings.HasPrefix(cmd.Use, "doctor") {
		initKubeClient()
	}

	cluster.SetImage(viper.GetString("image"))

	cluster.SetErrorHandlers()
}

func initKubeClient() {
	// The act of testing for a config, initializes the config.
	if err := doctor.IsClusterConfigValid(); err != nil {
		log.Fatal(err)
	}
}
