package ksync

import (
	"testing"

	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func init() {
// 	InitKubeClient("", "kube-system")
// }

// TODO: Change the function name to prevent collisions?
func TestMirrorRun(t *testing.T) {
	mirror := &Mirror{
		RemoteContainer: &RemoteContainer{
			// TODO: This has to be dynamic
			// See https://github.com/vapor-ware/ksync/blob/testier/pkg/ksync/container_test.go#L20
			NodeName: "gke-tim-dev-default-pool-9e45a876-pzbw",
		},
		// TODO: Need to make sure this always exists
		LocalPath: "/tmp/test",
		// TODO: Ditto
		RemotePath: "/tmp/test",
	}

	err := mirror.Run()

	require.NoError(t, err)
}
