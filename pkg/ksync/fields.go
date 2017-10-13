package ksync

import (
	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
)

func MergeFields(fieldSlice ...log.Fields) log.Fields {
	fields := &log.Fields{}
	for _, src := range fieldSlice {
		mergo.Merge(fields, src)
	}

	return *fields
}
