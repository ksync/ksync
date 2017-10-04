package main

import (
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

	// TODO: this is technically working like `find` right now. Should it be a
	// find or more like list?
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

	specList, err := ksync.AllSpecs()
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

	specList.Add(newSpec)
	if err := specList.Save(); err != nil {
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

	addViper.BindPFlag("name", runCmd.Flags().Lookup("name"))
	addViper.BindEnv("name", "KSYNC_NAME")
}
