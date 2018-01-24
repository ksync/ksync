package cluster

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	if os.Getenv("IMAGE") != "" {
		SetImage(os.Getenv("IMAGE"))
	}
}

func TestNewRadarInstance(t *testing.T) {
	service := NewService()

	require.NotPanics(t, func() { NewService() })
	assert.NotEmpty(t, service)
}

func TestServiceRun(t *testing.T) {
	service := NewService()
	require.NotPanics(t, func() { NewService() })

	// Normal run without upgrade
	err := service.Run(false)

	assert.NoError(t, err)
	assert.NotEmpty(t, service)

	// Run with upgrade
	// TODO: Use a new `radar` object here?
	err = service.Run(true)

	assert.NoError(t, err)
	assert.NotEmpty(t, service)
}
