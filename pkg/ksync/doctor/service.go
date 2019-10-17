package doctor

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/ksync/ksync/pkg/ksync"
	"github.com/ksync/ksync/pkg/ksync/cluster"
)

var (
	missingServiceError = `The cluster service has not been installed yet. Run init to fix.`
	serviceHealthError  = `Cluster service is not healthy.

- If you just ran init, wait a little longer and try again.
- Run 'kubectl --namespace=%s --context=%s get pods -lapp=ksync' to look at what's going on.`

	versionMismatch = `There is a mismatch between the local version (%s) and the cluster (%s).

Run 'ksync init --upgrade' to fix.`
)

// HasClusterService verifies that the cluster is running the service.
func HasClusterService() error {
	result, err := cluster.NewService().IsInstalled()
	if err != nil {
		return err
	}

	if !result {
		return fmt.Errorf(missingServiceError)
	}

	return nil
}

// IsClusterServiceHealthy verifies that the cluster service is healthy
// across all nodes.
func IsClusterServiceHealthy() error {
	// Note: this assumes that the service has already been added, otherwise the
	// error might not make as much sense. `HasClusterService` should have been
	// run beforehand.
	s := cluster.NewService()

	unhealthyError := fmt.Errorf(
		serviceHealthError,
		cluster.NewService().Namespace,
		viper.GetString("context"))

	nodes, err := s.NodeNames()
	if err != nil {
		return err
	} else if len(nodes) == 0 {
		return unhealthyError
	}

	for _, node := range nodes {
		if state, healthErr := s.IsHealthy(node); healthErr != nil {
			return healthErr
		} else if !state {
			return unhealthyError
		}
	}

	return nil
}

// IsServiceCompatible verifies that the remote service is compatible with
// the local client.
func IsServiceCompatible() error {
	version, err := cluster.NewService().Version()
	if err != nil {
		return err
	}

	if version.GitTag != ksync.GitTag {
		return fmt.Errorf(versionMismatch, ksync.GitTag, version.GitTag)
	}

	return nil
}
