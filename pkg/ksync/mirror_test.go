package ksync

// import (
//  "testing"
//  "os"
//
// 	// "github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func init() {
// 	InitKubeClient("", os.Getenv("TEST_NAMESPACE"))
// }

// TODO: Change the function name to prevent collisions?
// func TestMirrorRun(t *testing.T) {
// 	mirror := &Mirror{
// 		RemoteContainer: &RemoteContainer{
// 			// TODO: This has to be dynamic
// 			// See https://github.com/vapor-ware/ksync/blob/testier/pkg/ksync/container_test.go#L20
// 			NodeName: "gke-tim-dev-default-pool-9e45a876-2ggc",
// 			ID: "98ab1831587a84052db5f823a1a5742af2045c11b2f59a68ccf5f86ceb37a93f",
// 		},
// 		// TODO: Need to make sure this always exists
// 		LocalPath: "/tmp/test",
// 		// TODO: Ditto
// 		RemotePath: "/tmp/test",
// 	}
//
// 	err := mirror.Run()
//
// 	require.NoError(t, err)
// }
