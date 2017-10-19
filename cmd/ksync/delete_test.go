package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

func TestDeleteNew(t *testing.T) {
	testCobra := &deleteCmd{}
	cmd := testCobra.new()

	assert.IsTypef(t, reflect.TypeOf(&cobra.Command{}), reflect.TypeOf(cmd), "New command is of type %s", reflect.TypeOf(cmd))
	// TODO: Write more specific test cases
}
