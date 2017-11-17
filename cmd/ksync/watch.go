package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
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

// TODO: hook up to k8s and watch for changes
// TODO: handle Normalize errors.
func (w *watchCmd) run(cmd *cobra.Command, args []string) {
	if !ksync.HasMirror() {
		log.Fatal("missing required files. run `ksync init` again.")
	}

	list := &ksync.SpecList{}
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

	waitForSignal()
}

func waitForSignal() {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM) // nolint: megacheck
	<-exitSignal
}
