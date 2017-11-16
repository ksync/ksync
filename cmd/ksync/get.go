package main

import (
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
)

type getCmd struct {
	cli.BaseCmd
}

func (g *getCmd) new() *cobra.Command {
	long := `Get all configured syncs and their status.`
	example := ``

	g.Init("ksync", &cobra.Command{
		Use:     "get",
		Short:   "Get all configured syncs and their status.",
		Long:    long,
		Example: example,
		Run:     g.run,
	})

	return g.Cmd
}

// TODO: add last_sync (last_run?)
// TODO: make the columns configurable
// TODO: add a quiet ouput that can be `ksync get -q | ksync delete`
// TODO: output different formats (json)
// TODO: make output configurable
// TODO: the paths can be pretty long, keep them to a certain length?
// TODO: check for existence of the watcher, warn if it isn't running.
func (g *getCmd) run(cmd *cobra.Command, args []string) {
	specs := &ksync.SpecList{}
	if err := specs.Update(); err != nil {
		log.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetColumnSeparator(" ")
	table.SetHeader([]string{"Name", "Local", "Remote", "Status"})

	var keys []string
	for name := range specs.Items {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	for _, name := range keys {
		spec := specs.Items[name]
		status, err := spec.Status()
		if err != nil {
			log.Fatal(err)
		}

		table.Append([]string{
			name,
			spec.LocalPath,
			spec.RemotePath,
			status,
		})
	}

	table.Render()
}
