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

var (
	// TODO: update the usage instructions
	addHelp = `
    Add a new sync between a local and remote directory.
    `

	addCmd = &cobra.Command{
		Use:     "add [flags] [local path] [remote path]",
		Short:   "Add a new sync between a local and remote directory.",
		Long:    addHelp,
		Aliases: []string{"a"},
		Args:    cobra.ExactArgs(2),
		Run:     runAdd,
		// TODO: BashCompletionFunction
	}

	addViper = viper.New()
)

func runAdd(_ *cobra.Command, args []string) {
	loc := input.GetLocator(addViper)
	paths := input.GetPaths(args)

	// Usage validation ------------------------------------
	loc.Validator()
	paths.Validator()

	name := addViper.GetString("name")
	if name == "" {
		rand.Seed(time.Now().UnixNano())
		name = petname.Generate(2, "-")
	}

	specMap, err := ksync.AllSpecs()
	if err != nil {
		log.Fatal(err)
	}

	newSpec := &ksync.Spec{
		Container:  addViper.GetString("container"),
		Pod:        addViper.GetString("pod"),
		Selector:   addViper.GetString("selector"),
		LocalPath:  paths.Local,
		RemotePath: paths.Remote,
	}

	if err := specMap.Add(name, newSpec, addViper.GetBool("force")); err != nil {
		log.Fatalf("Could not add, --force to ignore: %v", err)
	}
	if err := specMap.Save(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	RootCmd.AddCommand(addCmd)

	input.LocatorFlags(addCmd, addViper)

	addCmd.Flags().String(
		"name",
		"",
		"Friendly name to describe this sync.")

	addViper.BindPFlag("name", addCmd.Flags().Lookup("name"))
	addViper.BindEnv("name", "KSYNC_NAME")

	addCmd.Flags().Bool(
		"force",
		false,
		"Force addition, ignoring similarity.")

	addViper.BindPFlag("force", addCmd.Flags().Lookup("force"))
	addViper.BindEnv("force", "KSYNC_FORCE")
}
