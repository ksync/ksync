package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync/doctor"
)

var (
	doctorSuccess = `Everything looks good!`
)

type doctorCmd struct {
	cli.BaseCmd
}

func (d *doctorCmd) new() *cobra.Command {
	long := `Troubleshoot and verify your setup is correct.`
	example := ``

	d.Init("ksync", &cobra.Command{
		Use:     "doctor [flags] [name]",
		Short:   "Troubleshoot and verify your setup is correct.",
		Long:    long,
		Example: example,
		Run:     d.run,
	})

	flags := d.Cmd.Flags()

	flags.Bool(
		"ignore",
		false,
		"Ignore test failures and try all the checks.")
	if err := d.BindFlag("ignore"); err != nil {
		log.Fatal(err)
	}

	return d.Cmd
}

func (d *doctorCmd) run(cmd *cobra.Command, args []string) {
	failure := false

	for _, check := range doctor.CheckList {
		if err := check.Out(); err != nil {
			if !d.Viper.GetBool("ignore") {
				failure = true
				break
			}
		}
	}

	if !failure {
		fmt.Printf(doctorSuccess)
	}
}
