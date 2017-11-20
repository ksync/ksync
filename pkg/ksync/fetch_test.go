package ksync

import (
  "testing"

  "github.com/stretchr/testify/assert"
  // "github.com/stretchr/testify/require"
)

func TestFetchMirror(t *testing.T) {
  err := FetchMirror()

  assert.NoError(t, err)
}
