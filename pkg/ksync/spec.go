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

// SpecMap is a collection of Specs.
type SpecMap struct {
	Items map[string]*Spec
}

// Spec is all the configuration required to setup a sync from a local directory
// to a remote directory in a specific remote container.
type Spec struct {
	Container string
	// TODO: use a locator instead?
	Pod        string
	Selector   string
	LocalPath  string
	RemotePath string
}

func (this *SpecMap) String() string {
	return YamlString(this)
}

func (this *SpecMap) Fields() log.Fields {
	return log.Fields{}
}

func (this *Spec) String() string {
	return YamlString(this)
}

func (this *Spec) Fields() log.Fields {
	return StructFields(this)
}

// AllSpecs populates a SpecMap with the configured specs. These are populated
// normally via. configuration.
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

// Create checks an individual input spec for likeness and duplicates
// then adds the spec into a SpecMap
func (this *SpecMap) Create(name string, spec *Spec, force bool) error {
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

// Delete removes a given spec from a SpecMap
func (this *SpecMap) Delete(name string) error {
	if !this.Has(name) {
		return fmt.Errorf("does not exist")
	}

	delete(this.Items, name)
	return nil
}

// Save serializes the current SpecMap's items to the config file.
func (this *SpecMap) Save() error {
	cfgPath := viper.ConfigFileUsed()
	if cfgPath == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		cfgPath = filepath.Join(home, fmt.Sprintf(".%s.yaml", "ksync"))
	}

	fobj, err := os.OpenFile(cfgPath, os.O_CREATE|os.O_WRONLY, 0644)
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

// HasLike checks a given spec for deep equivalence against another spec
// TODO: is this the best way to do this?
func (this *SpecMap) HasLike(target *Spec) bool {
	for _, spec := range this.Items {
		if reflect.DeepEqual(target, spec) {
			return true
		}
	}
	return false
}

// Has checks a given spec for simple equivalence against another spec
func (this *SpecMap) Has(target string) bool {
	if _, ok := this.Items[target]; ok {
		return true
	}
	return false
}
