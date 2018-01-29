package main

import (
	"github.com/spf13/cobra"

	"github.com/jpillora/overseer/fetcher"
	"github.com/timfallmk/overseer"

	"github.com/vapor-ware/ksync/pkg/cli"
)

const (
	repoUsername = "vapor-ware"
	repoName     = "ksync"
)

type updateCmd struct {
	cli.BaseCmd
}

func (u *updateCmd) new() *cobra.Command {
	long := `Update ksync to the latest version.`
	example := ``

	u.Init("ksync", &cobra.Command{
		Use:     "update",
		Short:   "Update ksync to the latest version.",
		Long:    long,
		Example: example,
		Run:     u.run,
	})

	return u.Cmd
}

func (u *updateCmd) run(cmd *cobra.Command, args []string) {
	overseer.Run(overseer.Config{
		Required:  true,
		Program:   func(_ overseer.State) {},
		Address:   ":0000",
		NoRestart: true,
		Debug:     false,
		Fetcher: &fetcher.Github{
			User: repoUsername,
			Repo: repoName,
		},
	})
}
