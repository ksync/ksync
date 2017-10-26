package ksync

import (
	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
)

// MergeFields takes a slice of logging fields and merges them together.
func MergeFields(fieldSlice ...log.Fields) (log.Fields, error) {
	fields := &log.Fields{}
	for _, src := range fieldSlice {
		err := mergo.Merge(fields, src)
		if err != nil {
			return nil, err
		}
	}

	return *fields, nil
}
