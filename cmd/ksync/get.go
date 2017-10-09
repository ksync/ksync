package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type GetCmd struct{}

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

func (this *GetCmd) run(cmd *cobra.Command, args []string) {
	log.Print("get cmd")
}
