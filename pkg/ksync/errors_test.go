package ksync

import (
	"testing"
	"fmt"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"

)

var (
	serviceerror = &serviceRunningError{}
)

func TestError(t *testing.T) {
	errmessage := serviceerror.Error()

	assert.EqualError(t, fmt.Errorf(errmessage), "Error: Already running: *ksync.Service\n----------\nnull\n----------")
}

func TestIsServiceRunning(t *testing.T) {
	errmessage := serviceerror.Error()
	ok := IsServiceRunning(fmt.Errorf(errmessage))

	assert.False(t, ok)
}
