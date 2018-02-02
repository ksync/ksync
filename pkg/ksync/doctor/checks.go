package doctor

import (
	"github.com/vapor-ware/ksync/pkg/cli"
)

type Check struct {
	Name string
	Func func() error
	Type string
}

var CheckList = []Check{
	Check{
		Name: "Extra Binaries",
		Func: DoesSyncthingExist,
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
		Name: "Watch Running",
		Func: IsWatchRunning,
	},
}

func (c *Check) Out() error {
	return cli.TaskOut(c.Name, c.Func)
}
