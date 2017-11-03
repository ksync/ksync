package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
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

	i.Cmd.Flags().BoolP(
		"upgrade",
		"u",
		false,
		"Upgrade the currently running version.")
	if err := i.BindFlag("upgrade"); err != nil {
		log.Fatal(err)
	}

	i.Cmd.Flags().Bool(
		"force",
		false,
		"Force the upgrade to occur.")
	if err := i.BindFlag("force"); err != nil {
		log.Fatal(err)
	}

	i.Cmd.Flags().Bool(
		"watch",
		true,
		"Don't run watch in the background.")
	if err := i.BindFlag("watch"); err != nil {
		log.Fatal(err)
	}

	return i.Cmd
}

// TODO: add instructions for watchman and limits (and detect them)
// TODO: need a better error with instructions on how to fix errors starting radar
func (i *initCmd) run(cmd *cobra.Command, args []string) {
	upgrade := i.Viper.GetBool("upgrade")
	if err := ksync.NewRadarInstance().Run(upgrade); err != nil {
		log.Fatalf("could not start radar: %v", err)
	}

	if i.Viper.GetBool("watch") {
		if err := ksync.BackgroundWatch(
			viper.ConfigFileUsed(), upgrade); err != nil {

			log.Fatalf("could not start watch: %v", err)
		}
	}
}
