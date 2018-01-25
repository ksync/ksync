package main

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
	"github.com/vapor-ware/ksync/pkg/ksync/server"
)

type watchCmd struct {
	cli.BaseCmd
}

func (w *watchCmd) new() *cobra.Command {
	long := `Watch configured syncs and start them when required.

	Note: this is run automatically for you by 'ksync init'. It expects to be run
	from inside a container.`
	example := ``

	w.Init("ksync", &cobra.Command{
		Use:     "watch",
		Short:   "Watch configured syncs and start them when required.",
		Long:    long,
		Example: example,
		Run:     w.run,
	})

	flags := w.Cmd.Flags()
	flags.String(
		"bind",
		"127.0.0.1",
		"interface to which the server will bind")

	if err := w.BindFlag("bind"); err != nil {
		log.Fatal(err)
	}

	flags.BoolP(
		"daemon",
		"d",
		false,
		"Run the watch command in the background.")
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
		log.Debug("config change")
		if err := w.update(list); err != nil {
			log.Fatal(err)
		}
	})
}

// TODO: should the listen be random?
// TODO: does this need TLS?
func (w *watchCmd) run(cmd *cobra.Command, args []string) {
	list := &ksync.SpecList{}

	w.local(list)

	daemonize := w.Viper.GetBool("daemon")

	if daemonize {
		if err := ksync.NewSyncthing().Daemonize(); err != nil {
			log.Fatal(err)
		}
	} else if err := ksync.NewSyncthing().Run(); err != nil {
		log.Fatal(err)
	}

	if !daemonize {
		if err := server.Listen(
			list, w.Viper.GetString("bind"), viper.GetInt("port")); err != nil {
			log.Fatal(err)
		}
	}
}
