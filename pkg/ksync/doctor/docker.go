package doctor

import (
	"fmt"
	"reflect"

	"github.com/blang/semver"
	"github.com/golang/protobuf/ptypes/empty"
	// log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"

	"github.com/vapor-ware/ksync/pkg/ksync/cluster"
	pb "github.com/vapor-ware/ksync/pkg/proto"
)

var (
	dockerVersionError = `The docker version (%s) on node (%s) does not fall within the acceptible range for API versions: %s. Please upgrade to a compatible version.`
	dockerStorageError = `The configured docker storage driver (%s) on node (%s) is not part of the supported list: %s. Please open an issue to add support for your storage driver.`
	dockerGraphError   = `The configured docker storage root (%s) on node (%s) does not match the storage root specified: %s. Please check your remote storage root or pass the correct root in init with --graph-root.`
)

// IsDockerVersionCompatible verifies that the remote cluster is running a
// docker daemon with an API version that falls within the compatible range.
func IsDockerVersionCompatible() error {
	nodes, err := cluster.NewService().NodeNames()
	if err != nil {
		return err
	}

	versionRange, err := semver.ParseRange(DockerAPIRange)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		conn, err := cluster.NewConnection(node).Radar()
		if err != nil {
			return err
		}
		defer conn.Close() // nolint: errcheck

		info, err := pb.NewRadarClient(conn).GetDockerVersion(
			context.Background(), &empty.Empty{})
		if err != nil {
			return err
		}

		// Docker's API is not versioned like semver, so we make it that way for
		// fun and games.
		apiVersion, err := semver.Make(info.APIVersion + ".0")
		if err != nil {
			return err
		}

		if !versionRange(apiVersion) {
			return fmt.Errorf(
				dockerVersionError,
				info.Version,
				node,
				DockerRange)
		}
	}

	return nil
}

// IsDockerStorageCompatible verifies that the remote cluster has been
// configured to use compatible docker storage drivers.
func IsDockerStorageCompatible() error {
	nodes, err := cluster.NewService().NodeNames()
	if err != nil {
		return err
	}

	for _, node := range nodes {
		conn, err := cluster.NewConnection(node).Radar()
		if err != nil {
			return err
		}
		defer conn.Close() // nolint: errcheck

		info, err := pb.NewRadarClient(conn).GetDockerInfo(
			context.Background(), &empty.Empty{})
		if err != nil {
			return err
		}

		if _, ok := DockerDriver[info.Driver]; !ok {
			return fmt.Errorf(
				dockerStorageError,
				info.Driver,
				node,
				reflect.ValueOf(DockerDriver).MapKeys())
		}
	}

	return nil
}

// IsDockerGraphMatching checks to see if the configured graph directory (the location of the graph driver storage directories) matches what's configured in radar
func IsDockerGraphMatching() error {
	nodes, err := cluster.NewService().NodeNames()
	if err != nil {
		return err
	}

	for _, node := range nodes {
		conn, err := cluster.NewConnection(node).Radar()
		if err != nil {
			return err
		}
		defer conn.Close() // nolint: errcheck

		info, err := pb.NewRadarClient(conn).GetDockerInfo(
			context.Background(), &empty.Empty{})
		if err != nil {
			return err
		}

		if viper.GetString("graph-root") != info.GraphRoot {
			return fmt.Errorf(
				dockerGraphError,
				info.GraphRoot,
				node,
				viper.GetString("graph-root"))
		}
	}

	return nil
}
