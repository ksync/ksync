package doctor

import (
	"github.com/vapor-ware/ksync/pkg/cli"
)

// Check provides the definition of check that is named and validates something.
type Check struct {
	Name string
	Func func() error
	Type string
}

// CheckList is the full list of checks run by doctor.
var CheckList = []Check{
	Check{
		Name: "Extra Binaries",
		Func: DoesSyncthingExist,
	},
	Check{
		Name: "Local Config",
		Func: IsLocalConfigValid,
		Type: "pre",
	},
	Check{
		Name: "Cluster Config",
		Func: IsClusterConfigValid,
		Type: "pre",
	},
	Check{
		Name: "Cluster Connection",
		Func: CanConnectToCluster,
		Type: "pre",
	},
	Check{
		Name: "Cluster Version",
		Func: IsClusterVersionSupported,
		Type: "pre",
	},
	Check{
		Name: "Cluster Permissions",
		Func: HasClusterPermissions,
		Type: "pre",
	},
	Check{
		Name: "Cluster Service",
		Func: HasClusterService,
		Type: "post",
	},
	Check{
		Name: "Service Health",
		Func: IsClusterServiceHealthy,
		Type: "post",
	},
	Check{
		Name: "Service Version",
		Func: IsServiceCompatible,
		Type: "post",
	},
	Check{
		Name: "Docker Version",
		Func: IsDockerVersionCompatible,
		Type: "post",
	},
	Check{
		Name: "Docker Storage Driver",
		Func: IsDockerStorageCompatible,
		Type: "post",
	},
	Check{
		Name: "Docker Storage Root",
		Func: IsDockerGraphMatching,
		Type: "post",
	},
	Check{
		Name: "Watch Running",
		Func: IsWatchRunning,
	},
}

// Out provides pretty output with colors and spinners of progress.
func (c *Check) Out() error {
	return cli.TaskOut(c.Name, c.Func)
}
