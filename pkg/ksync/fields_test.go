package ksync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"
)

func TestMergeFields(t *testing.T) {
	fields := log.Fields{
		"field1": "1",
		"field2": "2",
		"field3": "3",
	}

	mergedfields, err := MergeFields(fields)

	require.NoError(t, err)
	assert.NotEmpty(t, mergedfields)
	assert.Equal(t, mergedfields, fields)
}
