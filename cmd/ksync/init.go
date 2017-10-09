package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/ksync"
)

type InitCmd struct{}

func (this *InitCmd) New() *cobra.Command {
	// TODO: update the usage instructions
	long := `
    Prepare the cluster.`
	example := ``

	cmd := &cobra.Command{
		Use:     "init [flags]",
		Short:   "Prepare the cluster.",
		Long:    long,
		Example: example,
		Run:     this.run,
	}

	flags := cmd.Flags()
	flags.BoolP(
		"upgrade",
		"u",
		false,
		"Upgrade the currently running version.")

	viper.BindPFlag("upgrade", flags.Lookup("upgrade"))
	viper.BindEnv("upgrade", "KSYNC_UPGRADE")

	flags.Bool(
		"force",
		false,
		"Force the upgrade to occur.")

	viper.BindPFlag("force", flags.Lookup("force"))
	viper.BindEnv("force", "KSYNC_FORCE")

	// TODO: client only flag

	return cmd
}

// TODO: add instructions for watchman and limits (and detect them)
func (this *InitCmd) run(cmd *cobra.Command, args []string) {
	err := ksync.InitRadar(viper.GetBool("upgrade"))
	// TODO: need a better error with instructions on how to fix it.
	if err != nil {
		log.Fatalf("could not start radar: %v", err)
	}
}
