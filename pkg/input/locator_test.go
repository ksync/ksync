package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vapor-ware/ksync/pkg/ksync"

)

var tl *Locator

func init(){
	tl = &Locator{
			PodName:        "testname",
			Selector:       "testselector",
			ContainerName:  "testcontainername",
		}
}

func TestValidator(t *testing.T) {
	assert.NotPanics(t, func() {tl.Validator()})
}

func TestContainers(t *testing.T) {
	cnerr := ksync.InitKubeClient("bluuuu", "bluuuuu")
	t.Log(tl)
	require.Error(t, cnerr)
	_, err := tl.Containers()

	assert.Error(t, err)
	tl = &Locator{PodName: "",}
}
