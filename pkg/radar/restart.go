package radar

import (
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	apiclient "github.com/docker/docker/client"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/spf13/viper"

	pb "github.com/ksync/ksync/pkg/proto"
)

// RestartSyncthing restarts the syncthing sidecar to this radar process. This is
// needed because:
//     When a container is started, it inherits the mount table. The mounts
//     are maintained internally. Any new mounts do not show up inside the
//     container. The implication of this is that any new containers starting
//     after syncthing will not have their FS mounted inside syncthing's container.
//     While the files are all available, the actual mount will not occur.
//     As syncthing clients will just reconnect after loosing connection with the
//     server, we restart syncthing to refresh the mounts on demand.
func (r *radarServer) RestartSyncthing(
	ctx context.Context, _ *empty.Empty) (*pb.Error, error) {

	// TODO: this is awful, I can't figure out how to attach config to context.
	podName := viper.GetString("pod-name")

	client, err := apiclient.NewClientWithOpts(apiclient.FromEnv)
	if err != nil {
		return nil, err
	}

	client.NegotiateAPIVersion(context.Background())

	args := filters.NewArgs()
	args.Add("label", "io.kubernetes.container.name=syncthing")
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
	}).Debug("found syncthing container")

	if err := client.ContainerRestart(
		context.Background(), cntr.ID, nil); err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"pod":    podName,
		"id":     cntr.ID,
		"status": cntr.Status,
		"state":  cntr.State,
	}).Debug("restarted syncthing container")

	return &pb.Error{Msg: ""}, nil
}

// Restart restarts a local container. This is an effective "hot reload" because
// docker restarts and keeps the overlayfs in place (we're still putting files
// into it)
func (r *radarServer) Restart(
	ctx context.Context, cntr *pb.ContainerPath) (*pb.Error, error) {

	client, err := apiclient.NewClientWithOpts(apiclient.FromEnv)
	if err != nil {
		return nil, err
	}

	client.NegotiateAPIVersion(context.Background())

	timeout := 0 * time.Second
	if err := client.ContainerRestart(
		context.Background(), cntr.ContainerId, &timeout); err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"id": cntr.ContainerId,
	}).Debug("restarted container")

	return &pb.Error{Msg: ""}, nil
}
