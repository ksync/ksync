package main

import (
	"strings"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	daemon "github.com/timfallmk/go-daemon"

	"github.com/ksync/ksync/pkg/cli"
	"github.com/ksync/ksync/pkg/ksync"
	"github.com/ksync/ksync/pkg/ksync/server"
)

type watchCmd struct {
	cli.BaseCmd
}

func (w *watchCmd) new() *cobra.Command {
	long := `Watch configured specs and start syncing files when required.`
	example := `ksync watch --daemon`

	w.Init("ksync", &cobra.Command{
		Use:     "watch",
		Short:   "Watch configured specs and start syncing files when required.",
		Long:    long,
		Example: example,
		Run:     w.run,
	})

	flags := w.Cmd.Flags()
	flags.String(
		"bind",
		"127.0.0.1",
		"interface to bind to")

	if err := w.BindFlag("bind"); err != nil {
		log.Fatal(err)
	}

	flags.BoolP(
		"daemon",
		"d",
		false,
		"run in the background")
	if err := w.BindFlag("daemon"); err != nil {
		log.Fatal(err)
	}

	return w.Cmd
}

func (w *watchCmd) update(list *ksync.SpecList) error {
	if err := list.Update(); err != nil {
		return err
	}

	if err := list.Watch(); err != nil {
		return err
	}

	return nil
}

func (w *watchCmd) local(list *ksync.SpecList) {
	if err := w.update(list); err != nil {
		log.Fatal(err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Ksync configuration change detected. Updating...")
		if err := w.update(list); err != nil {
			log.Fatal(err)
		}
	})
}

func (w *watchCmd) run(cmd *cobra.Command, args []string) {
	list := ksync.NewSpecList()

	w.local(list)

	if w.Viper.GetBool("daemon") {
		context := getDaemonContext()

		if _, err := context.Reborn(); err != nil {
			log.Fatal(err)
		}

		defer context.Release() //nolint: errcheck

		if !daemon.WasReborn() {
			log.Info("Sending watch to the background. Use clean to stop it.")
			return
		}
	}

	if err := ksync.NewSyncthing().Run(); err != nil {
		log.Fatal(err)
	}

	if err := server.Listen(
		list, w.Viper.GetString("bind"), viper.GetInt("port")); err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			log.Fatal("It appears that watch is already running. Stop it first.")
		}

		log.Fatal(err)
	}
}
