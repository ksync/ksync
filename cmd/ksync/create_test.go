package main

import (
  "testing"

  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

type MockCreateCmd struct {
  *CreateCmd
}

func init(create *CreateCmd) {
}

func TestNew(t *testing.T) {
  cmd := New()
  require.IsTypef(t, &Cobra.Command, cmd, "%s is not of type %s. Aborting.", cmd, Cobra.Command)

  assert.ObjectsAreEqualValues(t, MockCreateCmd, &cmd)
}
