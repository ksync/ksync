package main

import (
	"io/ioutil"
	"os"
	"syscall"

	daemon "github.com/sevlyar/go-daemon"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync/cluster"
)

type cleanCmd struct {
	cli.BaseCmd
}

func (c *cleanCmd) new() *cobra.Command {
	long := `Remove installed components.

WARNING: USING THE "NUKE" OPTION WILL REMOVE YOUR CONFIG. USE WITH CAUTION.
	`
	example := ``

	c.Init("ksync", &cobra.Command{
		Use:     "clean",
		Short:   "Remove installed pieces",
		Long:    long,
		Example: example,
		Run:     c.run,
		Hidden:  false,
	})

	flags := c.Cmd.Flags()
	flags.BoolP(
		"local",
		"l",
		false,
		"Remove local components (daemons, servers, etc.)")
	if err := c.BindFlag("local"); err != nil {
		log.Fatal(err)
	}

	flags.BoolP(
		"remote",
		"r",
		false,
		"Remove remote components (daemon-sets, pods, etc.)")
	if err := c.BindFlag("remote"); err != nil {
		log.Fatal(err)
	}

	flags.Bool(
		"nuke",
		false,
		"Remove everything including configs, db, and downloaded helper binaries. CAUTION!")
	if err := c.BindFlag("nuke"); err != nil {
		log.Fatal(err)
	}

	return c.Cmd
}

func (c *cleanCmd) cleanLocal() {
	context := getDaemonContext()
	child, err := context.Search()
	if err != nil {
		log.Infoln("No daemonized process found. Nothing to clean locally.")
		log.Warningln(err)
		return
	}

	// This is the dumbest thing in the world. We have to send signals using flags,
	// so create a new bool flag that is always true and passes SIGTERM.
	daemon.AddCommand(
		daemon.BoolFlag(func(b bool) *bool { return &b }(true)),
		syscall.SIGTERM,
		nil)

	// Skip error checking on this because this library is horrendous
	if err := daemon.SendCommands(child); err != nil {
		log.Fatal(err)
	}

	// Clean up after the process since it seems incapable of doing that itself
	if err := os.Remove(context.PidFileName); err != nil {
		log.Fatal(err)
	}
}

func (c *cleanCmd) cleanRemote() {
	service := cluster.NewService()

	// Check that the daemonset is running remotely
	if isInstalled, err := service.IsInstalled(); err != nil {
		log.Fatal(err)
	} else if !isInstalled {
		log.Infoln("Remote components are not installed")
		return
	}

	if err := service.Remove(); err != nil {
		log.Fatal(err)
	}
}

func (c *cleanCmd) fromOrbit() {
	log.Debug("Removing local processes")
	c.cleanLocal()
	log.Debug("Removing remote processes")
	c.cleanRemote()

	files, _ := ioutil.ReadDir(viper.ConfigFileUsed())
	log.WithFields(log.Fields{
		"path":  viper.ConfigFileUsed(),
		"files": files,
	}).Info("Nuking all files from from orbit. It's the only way.")
	if err := os.RemoveAll(cli.ConfigPath()); err != nil {
		log.Fatal(err)
	}
	log.Info("Nuke drop complete")
}

func (c *cleanCmd) run(cmd *cobra.Command, args []string) {
	if c.Viper.GetBool("local") {
		c.cleanLocal()
	}

	if c.Viper.GetBool("remote") {
		c.cleanRemote()
	}

	if c.Viper.GetBool("nuke") {
		c.fromOrbit()
	}

	if !c.Viper.GetBool("local") && !c.Viper.GetBool("remote") {
		c.cleanLocal()
		c.cleanRemote()
	}
}
