package cluster

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	if os.Getenv("IMAGE") != "" {
		SetImage(os.Getenv("IMAGE"))
	}

	// Set the default for `docker-root` so it's evaluated properly in the daemon set template during testing
	viper.Set("docker-root", "/var/lib/docker")
	// and same for docker-sock
	viper.Set("docker-socket", "/var/run/docker.sock")
	// and finally for the daemonset namespace
	viper.Set("daemonset-namespace", "kube-system")
}

func TestNewRadarInstance(t *testing.T) {
	service := NewService()

	require.NotPanics(t, func() { NewService() })
	assert.NotEmpty(t, service)
}

func TestServiceRun(t *testing.T) {
	service := NewService()
	require.NotPanics(t, func() { NewService() })

	// Normal run without upgrade & psp
	err := service.Run(false, false)

	assert.NoError(t, err)
	assert.NotEmpty(t, service)

	// Run with upgrade & psp
	// TODO: Use a new `radar` object here?
	err = service.Run(true, true)

	assert.NoError(t, err)
	assert.NotEmpty(t, service)
}
