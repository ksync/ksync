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
	createHelp = `
    create a new sync between a local and remote directory.
    `

	createCmd = &cobra.Command{
		Use:     "create [flags] [local path] [remote path]",
		Short:   "create a new sync between a local and remote directory.",
		Long:    createHelp,
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(2),
		Run:     runCreate,
		// TODO: BashCompletionFunction
	}

	createViper = viper.New()
)

func runCreate(_ *cobra.Command, args []string) {
	loc := input.GetLocator(createViper)
	paths := input.GetPaths(args)

	// Usage validation ------------------------------------
	loc.Validator()
	paths.Validator()

	name := createViper.GetString("name")
	if name == "" {
		rand.Seed(time.Now().UnixNano())
		name = petname.Generate(2, "-")
	}

	specMap, err := ksync.AllSpecs()
	if err != nil {
		log.Fatal(err)
	}

	newSpec := &ksync.Spec{
		Container:  createViper.GetString("container"),
		Pod:        createViper.GetString("pod"),
		Selector:   createViper.GetString("selector"),
		LocalPath:  paths.Local,
		RemotePath: paths.Remote,
	}

	if err := specMap.Create(name, newSpec, createViper.GetBool("force")); err != nil {
		log.Fatalf("Could not create, --force to ignore: %v", err)
	}
	if err := specMap.Save(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	RootCmd.AddCommand(createCmd)

	input.LocatorFlags(createCmd, createViper)

	createCmd.Flags().String(
		"name",
		"",
		"Friendly name to describe this sync.")

	createViper.BindPFlag("name", createCmd.Flags().Lookup("name"))
	createViper.BindEnv("name", "KSYNC_NAME")

	createCmd.Flags().Bool(
		"force",
		false,
		"Force createition, ignoring similarity.")

	createViper.BindPFlag("force", createCmd.Flags().Lookup("force"))
	createViper.BindEnv("force", "KSYNC_FORCE")
}
