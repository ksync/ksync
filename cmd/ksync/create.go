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

type CreateCmd struct {
	viper *viper.Viper
}

func (this *CreateCmd) New() *cobra.Command {
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
		Run:     this.run,
		// TODO: BashCompletionFunction
	}
	this.viper = viper.New()

	// TODO: can this become a mixin?
	input.LocatorFlags(cmd, this.viper)

	flags := cmd.Flags()
	flags.String(
		"name",
		"",
		"Friendly name to describe this sync.")

	this.viper.BindPFlag("name", flags.Lookup("name"))
	this.viper.BindEnv("name", "KSYNC_NAME")

	flags.Bool(
		"force",
		false,
		"Force createition, ignoring similarity.")

	this.viper.BindPFlag("force", flags.Lookup("force"))
	this.viper.BindEnv("force", "KSYNC_FORCE")

	return cmd
}

func (this *CreateCmd) run(cmd *cobra.Command, args []string) {
	loc := input.GetLocator(this.viper)
	paths := input.GetPaths(args)

	// Usage validation ------------------------------------
	loc.Validator()
	paths.Validator()

	name := this.viper.GetString("name")
	if name == "" {
		rand.Seed(time.Now().UnixNano())
		name = petname.Generate(2, "-")
	}

	specMap, err := ksync.AllSpecs()
	if err != nil {
		log.Fatal(err)
	}

	newSpec := &ksync.Spec{
		Container:  this.viper.GetString("container"),
		Pod:        this.viper.GetString("pod"),
		Selector:   this.viper.GetString("selector"),
		LocalPath:  paths.Local,
		RemotePath: paths.Remote,
	}

	if err := specMap.Create(name, newSpec, this.viper.GetBool("force")); err != nil {
		log.Fatalf("Could not create, --force to ignore: %v", err)
	}
	if err := specMap.Save(); err != nil {
		log.Fatal(err)
	}
}
