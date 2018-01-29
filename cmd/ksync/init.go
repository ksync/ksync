package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
	"github.com/vapor-ware/ksync/pkg/ksync/cluster"
)

type initCmd struct {
	cli.BaseCmd
}

// TODO: client only flag
func (i *initCmd) new() *cobra.Command {
	// TODO: update the usage instructions
	long := `
    Prepare the cluster.`
	example := ``

	i.Init("ksync", &cobra.Command{
		Use:     "init [flags]",
		Short:   "Prepare the cluster.",
		Long:    long,
		Example: example,
		Run:     i.run,
	})

	flags := i.Cmd.Flags()
	flags.BoolP(
		"upgrade",
		"u",
		false,
		"Upgrade the currently running version.")
	if err := i.BindFlag("upgrade"); err != nil {
		log.Fatal(err)
	}

	flags.Bool(
		"client",
		true,
		"Setup the client",
	)
	if err := i.BindFlag("client"); err != nil {
		log.Fatal(err)
	}

	flags.Bool(
		"server",
		true,
		"Setup the server",
	)
	if err := i.BindFlag("server"); err != nil {
		log.Fatal(err)
	}

	flags.Bool(
		"skip-checks",
		false,
		"Skip the environment validation checks.",
	)
	if err := i.BindFlag("skip-checks"); err != nil {
		log.Fatal(err)
	}

	return i.Cmd
}

func (i *initCmd) initServer() {
	upgrade := i.Viper.GetBool("upgrade")
	if err := cluster.NewService().Run(upgrade); err != nil {
		log.Fatalf("could not start radar: %v", err)
	}
}

func (i *initCmd) initClient() {
	if i.Viper.GetBool("skip-checks") {
		return
	}

	sync := ksync.NewSyncthing()
	if !sync.HasBinary() {
		if err := sync.Fetch(); err != nil {
			log.Fatal(err)
		}
	}
}

// TODO: add instructions for watchman and limits (and detect them)
// TODO: need a better error with instructions on how to fix errors starting radar
func (i *initCmd) run(cmd *cobra.Command, args []string) {
	if i.Viper.GetBool("server") {
		i.initServer()
	}

	if i.Viper.GetBool("client") {
		i.initClient()
	}
}
