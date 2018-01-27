package main

import (
	"path/filepath"
	"syscall"

	daemon "github.com/sevlyar/go-daemon"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/cli"
)

type cleanCmd struct {
	cli.BaseCmd
}

func (c *cleanCmd) new() *cobra.Command {
	long := `Remove installed components.`
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
		"Remove everything including configs, db, and downloaded helper binaries. WARNING: THIS WILL REMOVE YOUR CONFIG. USE WITH CAUTION.")
	if err := c.BindFlag("nuke"); err != nil {
		log.Fatal(err)
	}

	return c.Cmd
}

func (c *cleanCmd) cleanLocal() {
	rootDir := filepath.Dir(viper.ConfigFileUsed())
	context := &daemon.Context{
		PidFileName: filepath.Join(rootDir, "daemon.pid"),
		LogFileName: filepath.Join(rootDir, "daemon.log"),
		WorkDir:     rootDir}

	child, err := context.Search()
	if err != nil {
		log.Infoln("No daemonized process found. Nothing to clean locally.")
		log.Fatalln(err)
	}

	// This is the dumbest thing in the world. We have to send signals using flags,
	// so create a new bool flag that is always true and passes SIGTERM.
	daemon.AddCommand(
		daemon.BoolFlag(func(b bool) *bool { return &b }(true)),
		syscall.SIGTERM,
		nil)
	daemon.SendCommands(child)
}

func (c *cleanCmd) cleanRemote() {

}

func (c *cleanCmd) fromOrbit() {

}

func (c *cleanCmd) run(cmd *cobra.Command, args []string) {
	if c.Viper.GetBool("local") {
		c.cleanLocal()
	}
	if c.Viper.GetBool("remote") {
		c.cleanRemote()
	}
	if c.Viper.GetBool("nuke") {
		c.cleanLocal()
		c.cleanRemote()
		c.fromOrbit()
	} else {
		c.cleanLocal()
		c.cleanRemote()
	}
}
