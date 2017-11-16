package radar

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	apiclient "github.com/docker/docker/client"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/spf13/viper"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

// RestartMirror restarts the mirror sidecar to this radar process. This is
// needed because:
//     When a container is started, it inherits the mount table. The mounts
//     are maintained internally. Any new mounts do not show up inside the
//     container. The implication of this is that any new containers starting
//     after mirror will not have their FS mounted inside mirror's container.
//     While the files are all available, the actual mount will not occur.
//     As mirror clients will just reconnect after loosing connection with the
//     server, we restart mirror to refresh the mounts on demand.
func (r *radarServer) RestartMirror(
	ctx context.Context, _ *empty.Empty) (*pb.Error, error) {

	// TODO: this is awful, I can't figure out how to attach config to context.
	podName := viper.GetString("pod-name")

	client, err := apiclient.NewEnvClient()
	if err != nil {
		return nil, err
	}

	args := filters.NewArgs()
	args.Add("label", "io.kubernetes.container.name=mirror")
	args.Add("label", fmt.Sprintf("io.kubernetes.pod.name=%s", podName))

	cntrs, err := client.ContainerList(
		context.Background(),
		types.ContainerListOptions{
			Filters: args,
		},
	)

	if err != nil {
		return nil, err
	}

	if len(cntrs) == 0 {
		return nil, fmt.Errorf("could not find for pod: %s", podName)
	}

	cntr := cntrs[0]

	log.WithFields(log.Fields{
		"pod":    podName,
		"id":     cntr.ID,
		"status": cntr.Status,
		"state":  cntr.State,
	}).Debug("found mirror container")

	if err := client.ContainerRestart(
		context.Background(), cntr.ID, nil); err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"pod":    podName,
		"id":     cntr.ID,
		"status": cntr.Status,
		"state":  cntr.State,
	}).Debug("restarted mirror container")

	return &pb.Error{Msg: ""}, nil
}

// Restart restarts a local container. This is an effective "hot reload" because
// docker restarts:
//   - keep the overlayfs in place (we're still putting files into it)
func (r *radarServer) Restart(
	ctx context.Context, cntr *pb.ContainerPath) (*pb.Error, error) {

	client, err := apiclient.NewEnvClient()
	if err != nil {
		return nil, err
	}

	if err := client.ContainerRestart(
		context.Background(), cntr.ContainerId, nil); err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"id": cntr.ContainerId,
	}).Debug("restarted container")

	return &pb.Error{Msg: ""}, nil
}
