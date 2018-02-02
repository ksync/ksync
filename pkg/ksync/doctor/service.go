package doctor

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/ksync"
	"github.com/vapor-ware/ksync/pkg/ksync/cluster"
)

var (
	serviceHealthError = `Cluster service is not healthy.

- If you just ran init, wait a little longer and try again.
- Run 'kubectl --namespace=%s --context=%s get pods -lapp=ksync' to look at what's going on.`

	versionMismatch = `There is a mismatch between the local version (%s) and the cluster (%s).

Run 'ksync init --upgrade' to fix.`
)

func HasClusterService() error {
	result, err := cluster.NewService().IsInstalled()
	if err != nil {
		return err
	}

	if !result {
		return fmt.Errorf(
			"The cluster service has not been installed yet. Run init to fix.")
	}

	return nil
}

func IsClusterServiceHealthy() error {
	// Note: this assumes that the service has already been added, otherwise the
	// error might not make as much sense. `HasClusterService` should have been
	// run beforehand.
	s := cluster.NewService()

	nodes, err := s.NodeNames()
	if err != nil {
		return err
	}

	for _, node := range nodes {
		if state, healthErr := s.IsHealthy(node); healthErr != nil {
			return healthErr
		} else if !state {
			return fmt.Errorf(
				serviceHealthError, cluster.NewService().Namespace, viper.GetString("context"))
		}
	}

	return nil
}

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
