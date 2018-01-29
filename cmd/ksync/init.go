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

func (i *initCmd) new() *cobra.Command {
	long := `Prepare ksync.

	Both the local host and remote cluster are initialized.`
	example := ``

	i.Init("ksync", &cobra.Command{
		Use:     "init [flags]",
		Short:   "Prepare ksync.",
		Long:    long,
		Example: example,
		Args:    cobra.ExactArgs(0),
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
		"local",
		true,
		"Initialize the local environment.",
	)
	if err := i.BindFlag("local"); err != nil {
		log.Fatal(err)
	}

	flags.Bool(
		"remote",
		true,
		"Initialize the remote environment.",
	)
	if err := i.BindFlag("remote"); err != nil {
		log.Fatal(err)
	}

	return i.Cmd
}

func (i *initCmd) initRemote() {
	upgrade := i.Viper.GetBool("upgrade")
	if err := cluster.NewService().Run(upgrade); err != nil {
		log.Fatalf("could not start radar: %v", err)
	}
}

func (i *initCmd) initLocal() {
	sync := ksync.NewSyncthing()
	if !sync.HasBinary() {
		if err := sync.Fetch(); err != nil {
			log.Fatal(err)
		}
	}
}

// TODO: need a better error with instructions on how to fix errors starting radar
func (i *initCmd) run(cmd *cobra.Command, args []string) {
	if i.Viper.GetBool("local") {
		i.initLocal()
	}

	if i.Viper.GetBool("remote") {
		i.initRemote()
	}
}
