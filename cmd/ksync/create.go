package main

import (
	"math/rand"
	"time"

	"github.com/dustinkirkland/golang-petname"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/input"
	"github.com/vapor-ware/ksync/pkg/ksync"
)

type createCmd struct {
	viper *viper.Viper
}

func (c *createCmd) new() *cobra.Command {
	long := `
    create a new sync between a local and remote directory.`
	example := ``

	cmd := &cobra.Command{
		Use:     "create [flags] [local path] [remote path]",
		Short:   "create a new sync between a local and remote directory.",
		Long:    long,
		Example: example,
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(2),
		Run:     c.run,
		// TODO: BashCompletionFunction
	}
	c.viper = viper.New()

	// TODO: can this become a mixin?
	input.LocatorFlags(cmd, c.viper)

	flags := cmd.Flags()
	flags.String(
		"name",
		"",
		"Friendly name to describe this sync.")

	c.viper.BindPFlag("name", flags.Lookup("name"))
	c.viper.BindEnv("name", "KSYNC_NAME")

	flags.Bool(
		"force",
		false,
		"Force creation, ignoring similarity.")

	c.viper.BindPFlag("force", flags.Lookup("force"))
	c.viper.BindEnv("force", "KSYNC_FORCE")

	return cmd
}

// TODO: check for existence of the watcher, warn if it isn't running.
func (c *createCmd) run(cmd *cobra.Command, args []string) {
	loc := input.GetLocator(c.viper)
	syncPath := input.GetSyncPath(args)

	// Usage validation ------------------------------------
	loc.Validator()
	syncPath.Validator()

	name := c.viper.GetString("name")
	if name == "" {
		rand.Seed(time.Now().UnixNano())
		name = petname.Generate(2, "-")
	}

	specMap, err := ksync.AllSpecs()
	if err != nil {
		log.Fatal(err)
	}

	newSpec := &ksync.Spec{
		Container:  c.viper.GetString("container"),
		Pod:        c.viper.GetString("pod"),
		Selector:   c.viper.GetString("selector"),
		LocalPath:  syncPath.Local,
		RemotePath: syncPath.Remote,
	}

	if err := specMap.Create(
		name, newSpec, c.viper.GetBool("force")); err != nil {
		log.Fatalf("Could not create, --force to ignore: %v", err)
	}
	if err := specMap.Save(); err != nil {
		log.Fatal(err)
	}
}
