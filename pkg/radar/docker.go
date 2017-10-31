package radar

import (
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

// TODO: needs to be able to reference volumes
// TODO: what to do about paths that include volumes? two syncs? they're different
// directories on the host itself. Maybe an alert for v1?
func getRootPath(containerPath *pb.ContainerPath) (string, error) {
	cli, err := client.NewEnvClient()
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
