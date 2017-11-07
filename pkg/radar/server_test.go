package radar

import (
  "testing"

  "github.com/stretchr/testify/assert"
  // "github.com/stretchr/testify/require"

  "google.golang.org/grpc"

)

func TestNewServer(t *testing.T) {
  server := NewServer()

  assert.NotPanics(t, func() { NewServer() })
  assert.NotEmpty(t, server)
  assert.IsType(t, &grpc.Server{}, server)
}
