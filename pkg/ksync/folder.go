package ksync

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/events"
	"github.com/syncthing/syncthing/lib/fs"
	"github.com/syncthing/syncthing/lib/protocol"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/vapor-ware/ksync/pkg/debug"
	pb "github.com/vapor-ware/ksync/pkg/proto"
	"github.com/vapor-ware/ksync/pkg/syncthing"
)

var tooSoonReset = 3 * time.Second

type Folder struct {
	SpecName        string
	RemoteContainer *RemoteContainer
	Reload          bool
	LocalPath       string
	RemotePath      string
	Status          ServiceStatus

	id string

	localServer    *syncthing.Server
	localDeviceID  protocol.DeviceID
	remoteServer   *syncthing.Server
	remoteDeviceID protocol.DeviceID

	restartContainer chan bool
	stop             chan bool
}

func NewFolder(service *Service) *Folder {
	return &Folder{
		SpecName:        service.SpecDetails.Name,
		RemoteContainer: service.RemoteContainer,
		Reload:          service.SpecDetails.Reload,
		LocalPath:       service.SpecDetails.LocalPath,
		RemotePath:      service.SpecDetails.RemotePath,

		id: fmt.Sprintf("%s-%s",
			service.SpecDetails.Name, service.RemoteContainer.PodName),
		stop: make(chan bool),
	}
}

func (f *Folder) String() string {
	return debug.YamlString(f)
}

func (f *Folder) initErrorHandler() {
	// Setup the k8s runtime to fail on unreturnable error (instead of looping).
	runtime.ErrorHandlers = append(runtime.ErrorHandlers, func(fromHandler error) {
		// We assume that this is always for a connection refused error, to log
		// and ignore it.

		log.WithFields(log.Fields{
			"node": f.RemoteContainer.NodeName,
			"pod":  f.RemoteContainer.PodName,
		}).Debug("lost connection to remote")

		return
	})
}

// Fields returns a set of structured fields for logging.
func (f *Folder) Fields() log.Fields {
	return debug.StructFields(f)
}

func (f *Folder) path() (string, error) {
	client, err := f.RemoteContainer.Radar()
	if err != nil {
		return "", err
	}

	path, err := client.GetBasePath(
		context.Background(), &pb.ContainerPath{
			ContainerId: f.RemoteContainer.ID,
		})
	if err != nil {
		return "", err
	}

	return filepath.Join(path.Full, f.RemotePath), nil
}

func (f *Folder) initServers(apiPort int32) error {
	localServer, err := syncthing.NewServer(
		fmt.Sprintf("localhost:%d", viper.GetInt("syncthing-port")),
		viper.GetString("apikey"))
	if err != nil {
		return err
	}

	f.localServer = localServer
	f.localDeviceID = f.localServer.ID

	remoteServer, err := syncthing.NewServer(
		fmt.Sprintf("localhost:%d", apiPort),
		viper.GetString("apikey"))
	if err != nil {
		return err
	}

	f.remoteServer = remoteServer
	f.remoteDeviceID = f.remoteServer.ID

	return nil
}

func (f *Folder) hotReload() error {
	f.restartContainer = make(chan bool)
	conn, err := NewRadarInstance().RadarConnection(
		f.RemoteContainer.NodeName)
	if err != nil {
		return err
	}

	client := pb.NewRadarClient(conn)

	// TODO: this is pretty naive, there are definite edge cases here where the
	// reload will happen but not actually get some files.
	go func() {
		defer conn.Close() // nolint: errcheck

		tooSoon := false
		for {
			select {
			case <-f.restartContainer:
				if tooSoon {
					continue
				}
				tooSoon = true

				log.WithFields(f.RemoteContainer.Fields()).Info("issuing reload")
				f.Status = ServiceReloading

				if _, err := client.Restart(
					context.Background(), &pb.ContainerPath{
						ContainerId: f.RemoteContainer.ID,
					}); err != nil {
					log.WithFields(f.RemoteContainer.Fields()).Error(err)
					continue
				}

				log.WithFields(f.RemoteContainer.Fields()).Info("reloaded")
				f.Status = ServiceWatching
			case <-time.After(tooSoonReset):
				tooSoon = false
			case <-f.stop:
				return
			}
		}
	}()

	return nil
}

func (f *Folder) watchEvents() error {
	stream, err := f.localServer.Events()
	if err != nil {
		return err
	}

	go func() {
		for ev := range stream {
			data := ev.Data.(map[string]interface{})

			if val, ok := data["folder"]; !ok || (ok && (f.id != val)) {
				continue
			}

			switch ev.Type {
			case events.FolderSummary:
				log.WithFields(f.RemoteContainer.Fields()).Info("updating")
				f.Status = ServiceUpdating
			case events.FolderCompletion:
				log.WithFields(f.RemoteContainer.Fields()).Info("update complete")
				f.Status = ServiceWatching

				if f.Reload {
					f.restartContainer <- true
				}
			}
		}
	}()

	return nil
}

func (f *Folder) setDevices(listenerPort int32) error {
	host, err := os.Hostname()
	if err != nil {
		return err
	}

	localDevice := config.NewDeviceConfiguration(f.localServer.ID, host)

	remoteDevice := config.NewDeviceConfiguration(
		f.remoteServer.ID, f.RemoteContainer.PodName)
	remoteDevice.Addresses = []string{
		fmt.Sprintf("tcp://127.0.0.1:%d", listenerPort),
	}

	if err := f.remoteServer.SetDevice(&localDevice); err != nil {
		return err
	}

	if err := f.localServer.SetDevice(&remoteDevice); err != nil {
		return err
	}

	return nil
}

func (f *Folder) setFolders() error {
	localFolder := config.NewFolderConfiguration(
		f.remoteServer.ID, f.id, f.id, fs.FilesystemTypeBasic, f.LocalPath)

	remotePath, err := f.path()
	if err != nil {
		return err
	}

	remoteFolder := config.NewFolderConfiguration(
		f.localServer.ID, f.id, f.id, fs.FilesystemTypeBasic, remotePath)

	if err := f.localServer.SetFolder(&localFolder); err != nil {
		return err
	}

	if err := f.remoteServer.SetFolder(&remoteFolder); err != nil {
		return err
	}

	return nil
}

func (f *Folder) Run() error {
	f.Status = ServiceStarting

	apiPort, listenerPort, err := NewRadarInstance().SyncthingConnection(
		f.RemoteContainer.NodeName)
	if err != nil {
		return err
	}
	f.initErrorHandler()

	if f.Reload {
		if err := f.hotReload(); err != nil {
			return err
		}
	}

	if err := f.initServers(apiPort); err != nil {
		return err
	}

	if err := f.watchEvents(); err != nil {
		return err
	}

	if err := f.setDevices(listenerPort); err != nil {
		return err
	}

	if err := f.setFolders(); err != nil {
		return err
	}

	if err := f.remoteServer.Update(); err != nil {
		return err
	}

	return f.localServer.Update()
}

// TODO: clear out remote config before leaving.
func (f *Folder) Stop() error {
	return nil
}
