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
	cli.FinderCmd
}

func (cmd *createCmd) new() *cobra.Command {
	long := `
    create a new sync between a local and remote directory.`
	example := ``

	cmd.Init("ksync", &cobra.Command{
		Use:     "create [flags] [local path] [remote path]",
		Short:   "create a new sync between a local and remote directory.",
		Long:    long,
		Example: example,
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(2),
		Run:     cmd.run,
		// TODO: BashCompletionFunction
	})

	if err := cmd.DefaultFlags(); err != nil {
		log.Fatal(err)
	}

	flags := cmd.Cmd.Flags()
	flags.String(
		"name",
		"",
		"Friendly name to describe this sync.")
	if err := cmd.BindFlag("name"); err != nil {
		log.Fatal(err)
	}

	flags.Bool(
		"force",
		false,
		"Force creation, ignoring similarity.")
	if err := cmd.BindFlag("force"); err != nil {
		log.Fatal(err)
	}

	return cmd.Cmd
}

// TODO: check for existence of the watcher, warn if it isn't running.
func (cmd *createCmd) run(_ *cobra.Command, args []string) {
	syncPath := input.GetSyncPath(args)

	// Usage validation ------------------------------------
	if err := cmd.Validator(); err != nil {
		log.Fatal(err)
	}
	syncPath.Validator()

	name := cmd.Viper.GetString("name")
	if name == "" {
		rand.Seed(time.Now().UnixNano())
		name = petname.Generate(2, "-")
	}

	specMap, err := ksync.AllSpecs()
	if err != nil {
		log.Fatal(err)
	}

	newSpec := &ksync.Spec{
		Container:  cmd.Viper.GetString("container"),
		Pod:        cmd.Viper.GetString("pod"),
		Selector:   cmd.Viper.GetString("selector"),
		LocalPath:  syncPath.Local,
		RemotePath: syncPath.Remote,
	}

	if err := specMap.Create(
		name, newSpec, cmd.Viper.GetBool("force")); err != nil {
		log.Fatalf("Could not create, --force to ignore: %v", err)
	}
	if err := specMap.Save(); err != nil {
		log.Fatal(err)
	}
}
