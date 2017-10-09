package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/input"
	"github.com/vapor-ware/ksync/pkg/ksync"
)

type ListCmd struct {
	viper *viper.Viper
}

func (this *ListCmd) New() *cobra.Command {
	long := `
    List the files from a remote container.`
	example := ``

	cmd := &cobra.Command{
		Use:     "list [flags] [path]",
		Short:   "List files from a remote container.",
		Long:    long,
		Example: example,
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		Run:     this.run,
		// TODO: BashCompletionFunction
	}
	this.viper = viper.New()

	// TODO: can this become a mixin?
	input.LocatorFlags(cmd, this.viper)

	return cmd
}

func (this *ListCmd) run(cmd *cobra.Command, args []string) {
	loc := input.GetLocator(this.viper)
	// Usage validation ------------------------------------
	loc.Validator()

	path := args[0]

	containerList, err := loc.Containers()
	if err != nil {
		log.Fatalf("%v", err)
	}

	// TODO: make this into a channel?
	// TODO: handle multi-container output
	for _, cntr := range containerList {
		list := &ksync.FileList{cntr, path, nil}
		if err := list.Get(); err != nil {
			log.Fatalf("%v", err)
		}

		if err := list.Output(); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
