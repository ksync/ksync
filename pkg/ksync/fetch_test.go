package ksync

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
)

func TestFetchMirror(t *testing.T) {
	err := FetchMirror()

	assert.NoError(t, err)

	// Cleanup
	// TODO: make the mirror jar variable
	cleanErr := os.Remove("mirror-all.jar")
	if os.IsNotExist(cleanErr) {
		t.Errorf("Binary was not found during cleanup. %s.", cleanErr)
	}
	t.Log("Test binary removed.")
}
