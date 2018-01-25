package ksync

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/fs"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/vapor-ware/ksync/pkg/debug"
	pb "github.com/vapor-ware/ksync/pkg/proto"
	"github.com/vapor-ware/ksync/pkg/syncthing"
)

type Folder struct {
	SpecName        string
	RemoteContainer *RemoteContainer
	Reload          bool
	LocalPath       string
	RemotePath      string
	Status          ServiceStatus

	local  *syncthing.Server
	remote *syncthing.Server
}

func NewFolder(service *Service) *Folder {
	return &Folder{
		SpecName:        service.SpecDetails.Name,
		RemoteContainer: service.RemoteContainer,
		Reload:          service.SpecDetails.Reload,
		LocalPath:       service.SpecDetails.LocalPath,
		RemotePath:      service.SpecDetails.RemotePath,
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
	local, err := syncthing.NewServer(
		fmt.Sprintf("localhost:%d", viper.GetInt("syncthing-port")),
		viper.GetString("apikey"))
	if err != nil {
		return err
	}

	f.local = local

	remote, err := syncthing.NewServer(
		fmt.Sprintf("localhost:%d", apiPort),
		viper.GetString("apikey"))
	if err != nil {
		return err
	}

	f.remote = remote

	return nil
}

func (f *Folder) setDevices(listenerPort int32) error {
	host, err := os.Hostname()
	if err != nil {
		return err
	}
	localDevice := config.NewDeviceConfiguration(f.local.ID, host)

	remoteDevice := config.NewDeviceConfiguration(
		f.remote.ID, f.RemoteContainer.PodName)
	remoteDevice.Addresses = []string{
		fmt.Sprintf("tcp://127.0.0.1:%d", listenerPort),
	}

	if err := f.remote.SetDevice(&localDevice); err != nil {
		return err
	}

	if err := f.local.SetDevice(&remoteDevice); err != nil {
		return err
	}

	return nil
}

func (f *Folder) setFolders() error {
	folderID := fmt.Sprintf("%s-%s", f.SpecName, f.RemoteContainer.PodName)
	localFolder := config.NewFolderConfiguration(
		f.remote.ID, folderID, folderID, fs.FilesystemTypeBasic, f.LocalPath)

	remotePath, err := f.path()
	if err != nil {
		return err
	}

	remoteFolder := config.NewFolderConfiguration(
		f.local.ID, folderID, folderID, fs.FilesystemTypeBasic, remotePath)

	if err := f.local.SetFolder(&localFolder); err != nil {
		return err
	}

	if err := f.remote.SetFolder(&remoteFolder); err != nil {
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

	if err := f.initServers(apiPort); err != nil {
		return err
	}

	if err := f.setDevices(listenerPort); err != nil {
		return err
	}

	if err := f.setFolders(); err != nil {
		return err
	}

	if err := f.remote.Update(); err != nil {
		return err
	}

	return f.local.Update()
}

func (f *Folder) Stop() error {
	return nil
}
