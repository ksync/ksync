package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
)

type deleteCmd struct {
	cli.BaseCmd
}

func (d *deleteCmd) new() *cobra.Command {
	long := `
		delete an existing sync.`
	example := ``

	d.Init("ksync", &cobra.Command{
		Use:     "delete [flags] [name]",
		Short:   "delete an existing sync.",
		Long:    long,
		Example: example,
		Aliases: []string{"d"},
		Args:    cobra.ExactArgs(1),
		Run:     d.run,
		// TODO: BashCompletionFunction
	})

	return d.Cmd
}

func (d *deleteCmd) run(cmd *cobra.Command, args []string) {
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
