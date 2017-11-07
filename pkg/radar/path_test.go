package radar

import (
	"testing"
	"context"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"

	pb "github.com/vapor-ware/ksync/pkg/proto"

)

func TestGetBasePath(t *testing.T) {
	radarserver := &radarServer{}
	context := context.Background()
	containerPath := &pb.ContainerPath{
    // TODO: This needs to be dynamic
    // TODO: This needs to reference a *local* docker container
		ContainerId: "",
	}
	basepath, err := (*radarserver).GetBasePath(context, containerPath)

	assert.Error(t, err)
	assert.Empty(t, basepath)
}
