package ksync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	servicetest = &Service{}
)

func init() {
  // TODO: Not sure this is strickly Kosher
	remotecontainer := &RemoteContainer{
		// TODO: This has to be dynamic
		// See https://github.com/vapor-ware/ksync/blob/testier/pkg/ksync/container_test.go#L20
		NodeName: "gke-tim-dev-default-pool-9e45a876-2ggc",
		ID: "98ab1831587a84052db5f823a1a5742af2045c11b2f59a68ccf5f86ceb37a93f",
	}
	spec := &Spec{}
	servicetest = NewService("test-service", remotecontainer, spec)
}
func TestGetServices(t *testing.T) {
	servicelist, err := GetServices()

	require.NotPanics(t, func() { GetServices() })

	assert.NoError(t, err)
	assert.NotEmpty(t, servicelist)
}
