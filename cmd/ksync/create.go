package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ksync/ksync/pkg/cli"
	"github.com/ksync/ksync/pkg/input"
	"github.com/ksync/ksync/pkg/ksync"
)

type createCmd struct {
	cli.FinderCmd
}

func (cmd *createCmd) new() *cobra.Command {
	long := `Create a new spec to sync files between a local and remote directory
  for specific containers running on the cluster.`
	example := `ksync create --local-read-only /code /go/src/github.com/ksync/code`

	cmd.Init("ksync", &cobra.Command{
		Use:     "create [flags] [local path] [remote path]",
		Short:   "Create a new spec",
		Long:    long,
		Example: example,
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(2),
		Run:     cmd.run,
	})

	if err := cmd.DefaultFlags(); err != nil {
		log.Fatal(err)
	}

	flags := cmd.Cmd.Flags()

	rand.Seed(time.Now().UnixNano())
	flags.String(
		"name",
		petname.Generate(2, "-"),
		"Friendly name to describe this spec.")
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

	flags.Bool(
		"reload",
		true,
		"Reload the remote container on file update.")
	if err := cmd.BindFlag("reload"); err != nil {
		log.Fatal(err)
	}

	flags.Bool(
		"local-read-only",
		false,
		"Set the local folder to read-only.")
	if err := cmd.BindFlag("local-read-only"); err != nil {
		log.Fatal(err)
	}

	flags.Bool(
		"remote-read-only",
		false,
		"Set the remote folder to read-only.")
	if err := cmd.BindFlag("remote-read-only"); err != nil {
		log.Fatal(err)
	}

	return cmd.Cmd
}

func (cmd *createCmd) run(_ *cobra.Command, args []string) {
	syncPath := input.GetSyncPath(args)

	// Usage validation ------------------------------------
	if err := cmd.Validator(); err != nil {
		log.Fatal(err)
	}
	if err := syncPath.Validator(); err != nil {
		if os.IsNotExist(err) {
			log.Fatal(fmt.Sprintf("local directory must exist (%s)", syncPath.Local))
		}
		log.Fatal(err)
	}

	specs := ksync.NewSpecList()
	if err := specs.Update(); err != nil {
		log.Fatal(err)
	}

	newSpec := &ksync.SpecDetails{
		Name: cmd.Viper.GetString("name"),

		ContainerName: cmd.Viper.GetString("container"),
		Pod:           cmd.Viper.GetString("pod"),
		Selector:      cmd.Viper.GetStringSlice("selector"),
		Namespace:     viper.GetString("namespace"),

		LocalReadOnly:  cmd.Viper.GetBool("local-read-only"),
		RemoteReadOnly: cmd.Viper.GetBool("remote-read-only"),

		LocalPath:  syncPath.Local,
		RemotePath: syncPath.Remote,

		Reload: cmd.Viper.GetBool("reload"),
	}

	if err := newSpec.IsValid(); err != nil {
		log.Fatal(err)
	}

	if err := specs.Create(
		newSpec,
		cmd.Viper.GetBool("force")); err != nil {

		log.Fatalf("Could not create, --force to ignore: %v", err)
	}

	if err := specs.Save(); err != nil {
		log.Fatal(err)
	}
}
