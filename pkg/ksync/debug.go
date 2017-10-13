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

func YamlString(thing interface{}) string {
	data, _ := yaml.Marshal(thing)
	return fmt.Sprintf(
		"%s\n----------\n%s----------", reflect.TypeOf(thing), string(data))
}

func StructFields(thing interface{}) log.Fields {
	return structs.Map(thing)
}

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
