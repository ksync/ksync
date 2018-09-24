package radar

import (
	"fmt"

	"github.com/docker/docker/client"
	apiclient "github.com/docker/docker/client"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

// TODO: needs to be able to reference volumes
// TODO: what to do about paths that include volumes? two syncs? they're different
// directories on the host itself. Maybe an alert for v1?
func getRootPath(containerPath *pb.ContainerPath) (string, error) {
	cli, err := apiclient.NewClientWithOpts(apiclient.FromEnv)
	if err != nil {
		return "", err
	}

	log.Debug("docker client created")

	cntr, err := cli.ContainerInspect(
		context.Background(), containerPath.ContainerId)
	if err != nil {
		return "", err
	}

	log.WithFields(log.Fields{
		"name": cntr.Name,
		"id":   containerPath.ContainerId,
	}).Debug("merge path retrieved")

	// TODO: how does this work on systems not running overlay2? Will need to
	// select on type.
	return cntr.GraphDriver.Data["MergedDir"], nil
}

func (r *radarServer) GetDockerVersion(
	ctx context.Context, _ *empty.Empty) (*pb.DockerVersion, error) {

	client, err := apiclient.NewClientWithOpts(apiclient.FromEnv)
	if err != nil {
		return nil, err
	}

	info, err := client.ServerVersion(context.Background())
	if err != nil {
		return nil, err
	}

	return &pb.DockerVersion{
		Version:       info.Version,
		APIVersion:    info.APIVersion,
		MinAPIVersion: info.MinAPIVersion,
		GitCommit:     info.GitCommit,
		GoVersion:     info.GoVersion,
		Os:            info.Os,
		Arch:          info.Arch,
	}, nil
}

func (r *radarServer) GetDockerInfo(
	ctx context.Context, _ *empty.Empty) (*pb.DockerInfo, error) {

	client, err := apiclient.NewClientWithOpts(apiclient.FromEnv)
	if err != nil {
		return nil, err
	}

	info, err := client.Info(context.Background())
	if err != nil {
		return nil, err
	}

	status := []string{}
	for _, pair := range info.DriverStatus {
		status = append(status, fmt.Sprintf("%s: %s", pair[0], pair[1]))
	}

	return &pb.DockerInfo{
		Driver:       info.Driver,
		DriverStatus: status,
		DockerRoot:   info.DockerRootDir,
	}, nil
}
