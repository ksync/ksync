package main

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
	"github.com/vapor-ware/ksync/pkg/ksync/cluster"
	"github.com/vapor-ware/ksync/pkg/ksync/doctor"
)

var (
	serviceHealthTimeout = `Timed out waiting for the ksync daemonset to be healthy.

To debug, you can:
- Run 'kubectl --namespace=%s --context=%s get pods -lapp=ksync' to look at what's going on.
- Run 'ksync doctor' to do an in-depth check of your system and the cluster.`
)

type initCmd struct {
	cli.BaseCmd
}

func (i *initCmd) new() *cobra.Command {
	long := `Prepare ksync.

	Both the local host and remote cluster are initialized.`
	example := `ksync init --local`

	i.Init("ksync", &cobra.Command{
		Use:     "init [flags]",
		Short:   "Prepare ksync.",
		Long:    long,
		Example: example,
		Args:    cobra.ExactArgs(0),
		Run:     i.run,
	})

	flags := i.Cmd.Flags()
	flags.BoolP(
		"upgrade",
		"u",
		false,
		"Upgrade the currently running version.")
	if err := i.BindFlag("upgrade"); err != nil {
		log.Fatal(err)
	}

	flags.Bool(
		"local",
		true,
		"Initialize the local environment.",
	)
	if err := i.BindFlag("local"); err != nil {
		log.Fatal(err)
	}

	flags.Bool(
		"remote",
		true,
		"Initialize the remote environment.",
	)
	if err := i.BindFlag("remote"); err != nil {
		log.Fatal(err)
	}

	flags.Bool(
		"skip-checks",
		false,
		"Skip the environment checks entirely.",
	)
	if err := i.BindFlag("skip-checks"); err != nil {
		log.Fatal(err)
	}

	return i.Cmd
}

func (i *initCmd) waitForHealthy() error {
	healthBackoff := backoff.NewExponentialBackOff()
	healthBackoff.MaxElapsedTime = 60 * time.Second

	service := cluster.NewService()
	if err := backoff.Retry(
		doctor.IsClusterServiceHealthy, healthBackoff); err != nil {

		log.Debug(err)
		return fmt.Errorf(
			serviceHealthTimeout, service.Namespace, viper.GetString("context"))
	}

	return nil
}

func (i *initCmd) remotePreChecks() {
	fmt.Println("==== Preflight checks ====")

	for _, check := range doctor.CheckList {
		if check.Type != "pre" {
			continue
		}

		if err := check.Out(); err != nil {
			log.Fatal("Fix errors and try again.")
		}
	}

	fmt.Println()
}

func (i *initCmd) remotePostChecks() {
	fmt.Println("==== Postflight checks ====")

	for _, check := range doctor.CheckList {
		if check.Type != "post" {
			continue
		}

		if err := check.Out(); err != nil {
			log.Fatal()
		}
	}

	fmt.Println()
}

func (i *initCmd) initRemote() {
	if !i.Viper.GetBool("skip-checks") {
		i.remotePreChecks()
	}

	fmt.Println("==== Cluster Environment ====")

	add := func() error {
		return cluster.NewService().Run(i.Viper.GetBool("upgrade"))
	}

	if err := cli.TaskOut("Adding ksync to the cluster", add); err != nil {
		log.Fatal()
	}

	if err := cli.TaskOut(
		"Waiting for pods to be healthy", i.waitForHealthy); err != nil {
		log.Fatal()
	}

	fmt.Println()

	if !i.Viper.GetBool("skip-checks") {
		i.remotePostChecks()
	}
}

func (i *initCmd) initLocal() {
	if err := doctor.DoesSyncthingExist(); err == nil {
		return
	}

	fmt.Println("==== Local Environment ====")

	if err := cli.TaskOut(
		"Fetching extra binaries", ksync.NewSyncthing().Fetch); err != nil {
		log.Fatal()
	}

	fmt.Println()
}

// TODO: need a better error with instructions on how to fix errors starting radar
func (i *initCmd) run(cmd *cobra.Command, args []string) {
	if i.Viper.GetBool("local") {
		i.initLocal()
	}

	if i.Viper.GetBool("remote") {
		i.initRemote()
	}

	fmt.Println("==== Initialization Complete ====")
}
