package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/input"
	"github.com/vapor-ware/ksync/pkg/ksync"
)

type listCmd struct {
	cli.BaseCmd
}

func (l *listCmd) new() *cobra.Command {
	long := `
    List the files from a remote container.`
	example := ``

	l.Init("ksync", &cobra.Command{
		Use:     "list [flags] [path]",
		Short:   "List files from a remote container.",
		Long:    long,
		Example: example,
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		Run:     l.run,
		// TODO: BashCompletionFunction
	})

	// TODO: can this become a mixin?
	input.LocatorFlags(l.Cmd, l.Viper)

	return l.Cmd
}

func (l *listCmd) run(cmd *cobra.Command, args []string) {
	loc := input.GetLocator(l.Viper)
	// Usage validation ------------------------------------
	loc.Validator()

	path := args[0]

	containerList, err := loc.Containers()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: make this into a channel?
	for _, cntr := range containerList {
		list := &ksync.FileList{
			Container: cntr,
			Path:      path,
		}
		if err := list.Get(); err != nil {
			log.Fatalf("%v", err)
		}

		if err := list.Output(); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
