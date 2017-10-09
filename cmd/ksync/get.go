package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vapor-ware/ksync/pkg/ksync"
)

// GetCmd specifies the structure of the `ksync get` command parameters
type GetCmd struct{}

// New creates a new `get` command and initializes the default values
func (this *GetCmd) New() *cobra.Command {
	long := ``
	example := ``

	cmd := &cobra.Command{
		Use:     "get",
		Short:   "",
		Long:    long,
		Example: example,
		Run:     this.run,
	}

	return cmd
}

// run takes the newly formed `get` command and combines it with general
// flags. These flags are then validated, before the entire command is run and
// any matching output displayed in a table.
// TODO: add last_sync (last_run?)
// TODO: make the columns configurable
// TODO: add a quiet ouput that can be `ksync get -q | ksync delete`
// TODO: output different formats (json)
// TODO: make output configurable
// TODO: the paths can be pretty long, keep them to a certain length?
func (this *GetCmd) run(cmd *cobra.Command, args []string) {
	specMap, err := ksync.AllSpecs()
	if err != nil {
		log.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetColumnSeparator(" ")
	table.SetHeader([]string{"Name", "Local", "Remote", "Status"})

	for name, spec := range specMap.Items {
		table.Append([]string{
			name,
			spec.LocalPath,
			spec.RemotePath,
			"Not Running",
		})
	}

	table.Render()
}
