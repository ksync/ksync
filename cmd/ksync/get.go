package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vapor-ware/ksync/pkg/ksync"
)

type getCmd struct{}

func (this *getCmd) new() *cobra.Command {
	long := `Get all configured syncs and their status.`
	example := ``

	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get all configured syncs and their status.",
		Long:    long,
		Example: example,
		Run:     this.run,
	}

	return cmd
}

// TODO: add last_sync (last_run?)
// TODO: make the columns configurable
// TODO: add a quiet ouput that can be `ksync get -q | ksync delete`
// TODO: output different formats (json)
// TODO: make output configurable
// TODO: the paths can be pretty long, keep them to a certain length?
// TODO: check for existence of the watcher, warn if it isn't running.
func (this *getCmd) run(cmd *cobra.Command, args []string) {
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
