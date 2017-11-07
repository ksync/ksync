package ksync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRadarInstance(t *testing.T) {
	radar := NewRadarInstance()

	require.NotPanics(t, func() { NewRadarInstance() })
	assert.NotEmpty(t, radar)
}

func TestRadarRun(t *testing.T) {
	radar := NewRadarInstance()
	require.NotPanics(t, func() { NewRadarInstance() })

	// Normal run without upgrade
	err := radar.Run(false)

	assert.NoError(t, err)
	assert.NotEmpty(t, radar)

	// Run with upgrade
	// TODO: Use a new `radar` object here?
	err = radar.Run(true)

	assert.NoError(t, err)
	assert.NotEmpty(t, radar)
}

func TestRadarConnection(t *testing.T) {
	// TODO: Have to figure out how to test this without a connection
}
