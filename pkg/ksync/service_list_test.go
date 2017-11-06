package ksync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetServices(t *testing.T) {
	service, err := GetServices()

	require.NotPanics(t, func() { GetServices() })

	assert.NoError(t, err)
	assert.NotEmpty(t, service)
}
