package ksync

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

type SpecList struct {
	Items []*Spec
}

type Spec struct {
	Container  string
	Pod        string
	Selector   string
	LocalPath  string
	RemotePath string
}

// TODO: test non-existant file
// TODO: test missing specs
func AllSpecs() (*SpecList, error) {
	var all SpecList

	if !viper.IsSet("spec") {
		return &all, nil
	}

	for _, raw := range viper.Get("spec").([]interface{}) {
		var spec Spec
		if err := mapstructure.Decode(raw, &spec); err != nil {
			return nil, errors.Wrap(err, "cannot get current specs")
		}

		all.Items = append(all.Items, &spec)
	}

	return &all, nil
}

func (this *SpecList) Add(spec *Spec) {
	if !this.Has(spec) {
		this.Items = append(this.Items, spec)
	}
}

func (this *SpecList) Save() error {
	cfgPath := viper.ConfigFileUsed()
	if cfgPath == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		cfgPath = filepath.Join(home, fmt.Sprintf(".%s.yaml", "ksync"))
	}

	fobj, err := os.Create(cfgPath)
	if err != nil {
		return err
	}

	defer fobj.Close()

	log.WithFields(log.Fields{
		"path": cfgPath,
	}).Debug("writing config file")

	buf, err := yaml.Marshal(viper.AllSettings())
	if err != nil {
		return err
	}

	if _, err := fobj.WriteString(string(buf)); err != nil {
		return err
	}

	return nil
}

// TODO: need some better equality logic here for ones that maybe need updating
// instead of addition/removal
func (this *SpecList) Has(target *Spec) bool {
	for _, current := range this.Items {
		if reflect.DeepEqual(current, target) {
			return true
		}
	}

	return false
}
