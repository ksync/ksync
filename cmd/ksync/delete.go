package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vapor-ware/ksync/pkg/ksync"
)

// DeleteCmd specifies the structure of the `ksync delete` command parameters
type DeleteCmd struct{}

// New creates a new `delete` command and initializes the default values
func (this *DeleteCmd) New() *cobra.Command {
	long := `
		delete an existing sync.`
	example := ``

	cmd := &cobra.Command{
		Use:     "delete [flags] [name]",
		Short:   "delete an existing sync.",
		Long:    long,
		Example: example,
		Aliases: []string{"d"},
		Args:    cobra.ExactArgs(1),
		Run:     this.run,
		// TODO: BashCompletionFunction
	}

	return cmd
}

// run takes the newly formed `delete` command and combines it with general
// flags. These flags are then validated, before the entire command is run to
// create a specification map.
func (this *DeleteCmd) run(cmd *cobra.Command, args []string) {
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
