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

type SpecMap struct {
	Items map[string]*Spec
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
func AllSpecs() (*SpecMap, error) {
	var all SpecMap
	all.Items = map[string]*Spec{}

	if !viper.IsSet("spec") {
		return &all, nil
	}

	for name, raw := range viper.GetStringMap("spec") {
		var spec Spec
		if err := mapstructure.Decode(raw, &spec); err != nil {
			return nil, errors.Wrap(err, "cannot get current specs")
		}

		all.Items[name] = &spec
	}

	return &all, nil
}

func (this *SpecMap) Add(name string, spec *Spec, force bool) error {
	if !force {
		if this.Has(name) {
			// TODO: make this into a type?
			return fmt.Errorf("name already exists.")
		}

		if this.HasLike(spec) {
			return fmt.Errorf("similar spec exists.")
		}
	}

	this.Items[name] = spec
	return nil
}

func (this *SpecMap) Save() error {
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

	viper.Set("spec", this.Items)
	buf, err := yaml.Marshal(viper.AllSettings())
	if err != nil {
		return err
	}

	if _, err := fobj.WriteString(string(buf)); err != nil {
		return err
	}

	return nil
}

// TODO: is this the best way to do this?
func (this *SpecMap) HasLike(target *Spec) bool {
	for _, spec := range this.Items {
		if reflect.DeepEqual(target, spec) {
			return true
		}
	}
	return false
}

func (this *SpecMap) Has(target string) bool {
	if _, ok := this.Items[target]; ok {
		return true
	}
	return false
}
