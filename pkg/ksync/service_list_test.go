package ksync

import (
	"testing"
	"os"

	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
)

var (
	servicetest = &Service{}
)

func init() {
  // TODO: Not sure this is strickly Kosher
	remotecontainer := &RemoteContainer{
		// TODO: This has to be dynamic
		// See https://github.com/vapor-ware/ksync/blob/testier/pkg/ksync/container_test.go#L20
		NodeName: os.Getenv("TEST_NODE"),
		ID: os.Getenv("TEST_CONTAINERID"),
	}
	spec := &Spec{}
	servicetest = NewService("test-service", remotecontainer, spec)
}
// TODO: Need to figure out how to do this outside of an actual client
// or the cluster
func TestGetServices(t *testing.T) {
	// servicelist, err := GetServices()
  //
	// require.NotPanics(t, func() { GetServices() })
  //
	// assert.NoError(t, err)
	// assert.NotEmpty(t, servicelist)
}
