package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vapor-ware/ksync/pkg/ksync"
)

var (
	// TODO: update the usage instructions
	removeHelp = `
    Remove an existing sync.
    `

	removeCmd = &cobra.Command{
		Use:     "remove [flags] [name]",
		Short:   "Remove an existing sync.",
		Long:    removeHelp,
		Aliases: []string{"r"},
		Args:    cobra.ExactArgs(1),
		Run:     runRemove,
		// TODO: BashCompletionFunction
	}
)

func runRemove(_ *cobra.Command, args []string) {
	name := args[0]

	specMap, err := ksync.AllSpecs()
	if err != nil {
		log.Fatal(err)
	}

	if !specMap.Has(name) {
		log.Fatalf("%s does not exist. Did you mean something else?", name)
	}

	if err := specMap.Remove(name); err != nil {
		log.Fatalf("Could not remove %s: %v", name, err)
	}
	if err := specMap.Save(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	RootCmd.AddCommand(removeCmd)
}
