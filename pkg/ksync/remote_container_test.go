package ksync

import (
	// "os"
	"testing"
	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"

	"github.com/vapor-ware/ksync/pkg/ksync/cluster"
)

func init() {
	cluster.InitKubeClient("") // nolint: errcheck
}

func TestRemoteContainer(t *testing.T) {
	// TODO: need tests
}
