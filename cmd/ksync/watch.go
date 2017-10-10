package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/ksync"
)

type WatchCmd struct {
	viper *viper.Viper
}

func (this *WatchCmd) New() *cobra.Command {
	long := `Watch configured syncs and start them when required.`
	example := ``

	cmd := &cobra.Command{
		Use:     "watch",
		Short:   "Watch configured syncs and start them when required.",
		Long:    long,
		Example: example,
		Run:     this.run,
	}

	this.viper = viper.New()

	return cmd
}

func (this *WatchCmd) run(cmd *cobra.Command, args []string) {
	// 1. Watch config file for updates
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		specMap, _ := ksync.AllSpecs()
		for name, spec := range specMap.Items {
			// Should run?
			containerList, err := ksync.GetContainers(
				spec.Pod, spec.Selector, spec.Container)
			if err != nil {
				log.Fatal(err)
			}

			for _, cntr := range containerList {
				// Is running?
				log.Print(cntr)
			}
			log.Print(name)
		}
	})

	// 2. Watch k8s API for updates
	// 3. Add/remove runs
	log.Print("watch cmd")

	waitForSignal()
}

func waitForSignal() {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}
