package debug

import (
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func TestMergeFields(t *testing.T) {
	fields := log.Fields{
		"field1": "1",
		"field2": "2",
		"field3": "3",
	}

	mergedfields := MergeFields(fields)

	assert.NotEmpty(t, mergedfields)
	assert.Equal(t, mergedfields, fields)
}
