package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
)

type listCmd struct {
	cli.FinderCmd
}

func (cmd *listCmd) new() *cobra.Command {
	long := `
    List the files from a remote container.`
	example := ``

	cmd.Init("ksync", &cobra.Command{
		Use:     "list [flags] [path]",
		Short:   "List files from a remote container.",
		Long:    long,
		Example: example,
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		Run:     cmd.run,
		// TODO: BashCompletionFunction
	})

	if err := cmd.DefaultFlags(); err != nil {
		log.Fatal(err)
	}

	return cmd.Cmd
}

func (cmd *listCmd) run(_ *cobra.Command, args []string) {
	// Usage validation ------------------------------------
	if err := cmd.Validator(); err != nil {
		log.Fatal(err)
	}

	path := args[0]

	containerList, err := cmd.Containers()
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
