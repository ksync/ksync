package ksync

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// YamlString is a YAML representation of an interface.
func YamlString(thing interface{}) string {
	data, _ := yaml.Marshal(thing)
	return fmt.Sprintf(
		"%s\n----------\n%s----------", reflect.TypeOf(thing), string(data))
}

// StructFields is the set of fields on an interface that get used for log
// output.
func StructFields(thing interface{}) log.Fields {
	return structs.Map(thing)
}

// ErrorOut is a convenience for constructing extremely informative debugging
// errors.
func ErrorOut(msg string, err error, thing interface{}) error {
	_, fn, line, _ := runtime.Caller(1)

	return errors.Wrap(
		err,
		fmt.Sprintf("msg: %s\nlocation: %s:%d\nstruct: %s\nnext",
			msg,
			fn,
			line,
			thing,
		),
	)
}
