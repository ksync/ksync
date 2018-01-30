package ksync

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/events"
	"github.com/syncthing/syncthing/lib/fs"

	"github.com/vapor-ware/ksync/pkg/debug"
	"github.com/vapor-ware/ksync/pkg/ksync/cluster"
	pb "github.com/vapor-ware/ksync/pkg/proto"
	"github.com/vapor-ware/ksync/pkg/syncthing"
)

var tooSoonReset = 3 * time.Second

// Folder is what controls the syncing between a local folder and a specific
// container running in the remote cluster.
type Folder struct {
	SpecName        string
	RemoteContainer *RemoteContainer
	Reload          bool
	LocalPath       string
	RemotePath      string
	Status          ServiceStatus

	id string

	localServer  *syncthing.Server
	remoteServer *syncthing.Server

	connection  *cluster.Connection
	radarConn   *grpc.ClientConn
	radarClient pb.RadarClient

	restartContainer chan bool
	stop             chan bool
}

// NewFolder constructs a Folder based off the provided Service.
func NewFolder(service *Service) *Folder {
	return &Folder{
		SpecName:        service.SpecDetails.Name,
		RemoteContainer: service.RemoteContainer,
		Reload:          service.SpecDetails.Reload,
		LocalPath:       service.SpecDetails.LocalPath,
		RemotePath:      service.SpecDetails.RemotePath,

		id: fmt.Sprintf("%s-%s",
			service.SpecDetails.Name, service.RemoteContainer.PodName),

		connection: cluster.NewConnection(service.RemoteContainer.NodeName),

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

// Get the remote folder's path from radar.
func (f *Folder) path() (string, error) {
	path, err := f.radarClient.GetBasePath(
		context.Background(), &pb.ContainerPath{
			ContainerId: f.RemoteContainer.ID,
		})
	if err != nil {
		return "", err
	}

	return filepath.Join(path.Full, f.RemotePath), nil
}

func (f *Folder) initRadarClient() error {
	conn, err := f.connection.Radar()
	if err != nil {
		return err
	}

	f.radarConn = conn
	f.radarClient = pb.NewRadarClient(conn)

	return nil
}

// Mounts are propogated for /var/lib/docker when a container is started. If
// the container in question was started after the syncthing container, the
// mount will not be present. This restarts just the syncthing container on the
// remote node to refresh the mount table. It will potentially interrupt any
// other folders operating on that host, but syncthing will quickly reconnect
// and startup again.
func (f *Folder) refreshSyncthing() error {
	if _, err := f.radarClient.RestartSyncthing(
		context.Background(), &empty.Empty{}); err != nil {
		return debug.ErrorLocation(err)
	}

	return nil
}

// Sets up syncthing.Server for both the local and remote server to allow for
// updating their configuration mutually.
func (f *Folder) initServers(apiPort int32) error {
	localServer, err := syncthing.NewServer(
		fmt.Sprintf("localhost:%d", viper.GetInt("syncthing-port")),
		viper.GetString("apikey"))
	if err != nil {
		return err
	}

	f.localServer = localServer

	remoteServer, err := syncthing.NewServer(
		fmt.Sprintf("localhost:%d", apiPort),
		viper.GetString("apikey"))
	if err != nil {
		return err
	}

	f.remoteServer = remoteServer

	return nil
}

// Kick the remote container when a folder has successfully completed updating.
// This is monitored from the local syncthing server.
func (f *Folder) hotReload() error {
	f.restartContainer = make(chan bool)

	// TODO: this is pretty naive, there are definite edge cases here where the
	// reload will happen but not actually get some files.
	go func() {
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

				if _, err := f.radarClient.Restart(
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

// Pay attention to the events coming off the local syncthing server to update
// state and reload the remote container if required.
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
		log.WithFields(f.Fields()).Debug("cleaning up event handler")
	}()

	return nil
}

// Update both the local and remote syncthing servers with devices allowing
// them to mutually connect (via. the local tunnel). None of the discovery
// or hole punching options in syncthing are used. The configuration forces
// connections *only* over the local tunnel.
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

// Update both the local and remote folder configuration for syncthing. Once
// this is updated, the syncing will actually start (assuming the devices can
// connect via. the local tunnel).
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

func (f *Folder) beginSync(listenerPort int32) error {
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

// Run starts syncing the folder between the local host and the remote
// container. It is expected that syncthing is already running locally (
// normally started by Syncthing).
func (f *Folder) Run() error {
	f.Status = ServiceStarting

	if err := f.initRadarClient(); err != nil {
		return err
	}

	if err := f.refreshSyncthing(); err != nil {
		return err
	}

	apiPort, listenerPort, err := f.connection.Syncthing()
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

	return f.beginSync(listenerPort)
}

// Stop cleans up everything running in the background. It removes the
// folder configuration from syncthing on the local/remote servers and
// destroys the tunnels to those servers (for this folder, as connections
// are per-folder).
func (f *Folder) Stop() error {
	f.localServer.Stop()
	f.remoteServer.Stop()

	close(f.stop)
	<-f.stop

	// Leave the devices, there might be other syncs with those nodes. It
	// shouldn't be a huge deal because the tunnel will be down unless active.
	f.localServer.RemoveFolder(f.id)
	f.remoteServer.RemoveFolder(f.id)

	if err := f.localServer.Update(); err != nil {
		return err
	}

	if err := f.remoteServer.Update(); err != nil {
		return err
	}

	f.radarConn.Close()
	f.connection.Stop()

	log.WithFields(f.Fields()).Debug("stopped folder")
	return nil
}
