package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vapor-ware/ksync/pkg/ksync"
)

var (
	// TODO: update the usage instructions
	deleteHelp = `
    delete an existing sync.
    `

	deleteCmd = &cobra.Command{
		Use:     "delete [flags] [name]",
		Short:   "delete an existing sync.",
		Long:    deleteHelp,
		Aliases: []string{"d"},
		Args:    cobra.ExactArgs(1),
		Run:     rundelete,
		// TODO: BashCompletionFunction
	}
)

func rundelete(_ *cobra.Command, args []string) {
	name := args[0]

	specMap, err := ksync.AllSpecs()
	if err != nil {
		log.Fatal(err)
	}

	if !specMap.Has(name) {
		log.Fatalf("%s does not exist. Did you mean something else?", name)
	}

	if err := specMap.Delete(name); err != nil {
		log.Fatalf("Could not delete %s: %v", name, err)
	}
	if err := specMap.Save(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
