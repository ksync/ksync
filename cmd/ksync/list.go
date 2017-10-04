package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/input"
	"github.com/vapor-ware/ksync/pkg/ksync"
)

var (
	// TODO: update the usage instructions
	listHelp = `
    List the files from a remote container.
    `

	// TODO: this is technically working like `find` right now. Should it be a
	// find or more like list?
	listCmd = &cobra.Command{
		Use:     "list [flags] [path]",
		Short:   "List files from a remote container.",
		Long:    listHelp,
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		Run:     runList,
		// TODO: BashCompletionFunction
	}

	listViper = viper.New()
)

func runList(_ *cobra.Command, args []string) {
	loc := input.GetLocator(listViper)
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

func init() {
	RootCmd.AddCommand(listCmd)

	input.LocatorFlags(listCmd, listViper)
}
