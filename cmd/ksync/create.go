package main

import (
	"math/rand"
	"time"

	"github.com/dustinkirkland/golang-petname"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/input"
	"github.com/vapor-ware/ksync/pkg/ksync"
)

type createCmd struct {
	cli.BaseCmd
}

func (c *createCmd) new() *cobra.Command {
	long := `
    create a new sync between a local and remote directory.`
	example := ``

	c.Init("ksync", &cobra.Command{
		Use:     "create [flags] [local path] [remote path]",
		Short:   "create a new sync between a local and remote directory.",
		Long:    long,
		Example: example,
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(2),
		Run:     c.run,
		// TODO: BashCompletionFunction
	})

	// TODO: can this become a mixin?
	input.LocatorFlags(c.Cmd, c.Viper)

	c.Cmd.Flags().String(
		"name",
		"",
		"Friendly name to describe this sync.")
	if err := c.BindFlag("name"); err != nil {
		log.Fatal(err)
	}

	c.Cmd.Flags().Bool(
		"force",
		false,
		"Force creation, ignoring similarity.")
	if err := c.BindFlag("force"); err != nil {
		log.Fatal(err)
	}

	return c.Cmd
}

// TODO: check for existence of the watcher, warn if it isn't running.
func (c *createCmd) run(cmd *cobra.Command, args []string) {
	loc := input.GetLocator(c.Viper)
	syncPath := input.GetSyncPath(args)

	// Usage validation ------------------------------------
	loc.Validator()
	syncPath.Validator()

	name := c.Viper.GetString("name")
	if name == "" {
		rand.Seed(time.Now().UnixNano())
		name = petname.Generate(2, "-")
	}

	specMap, err := ksync.AllSpecs()
	if err != nil {
		log.Fatal(err)
	}

	newSpec := &ksync.Spec{
		Container:  c.Viper.GetString("container"),
		Pod:        c.Viper.GetString("pod"),
		Selector:   c.Viper.GetString("selector"),
		LocalPath:  syncPath.Local,
		RemotePath: syncPath.Remote,
	}

	if err := specMap.Create(
		name, newSpec, c.Viper.GetBool("force")); err != nil {
		log.Fatalf("Could not create, --force to ignore: %v", err)
	}
	if err := specMap.Save(); err != nil {
		log.Fatal(err)
	}
}
