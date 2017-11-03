package ksync

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"

	"github.com/vapor-ware/ksync/pkg/debug"
)

// SpecMap is a collection of Specs.
type SpecMap struct {
	Items map[string]*Spec
}

// Spec is all the configuration required to setup a sync from a local directory
// to a remote directory in a specific remote container.
type Spec struct {
	// Local config
	Name string
	User string

	// Kubernetes config
	Namespace   string
	Context     string
	KubeCfgPath string

	// RemoteContainer Locator
	// TODO: use a locator instead?
	Container string
	Pod       string
	Selector  string

	// File config
	LocalPath  string
	RemotePath string
}

func (s *SpecMap) String() string {
	return debug.YamlString(s)
}

// Fields returns a set of structured fields for logging.
func (s *SpecMap) Fields() log.Fields {
	return log.Fields{}
}

func (s *Spec) String() string {
	return debug.YamlString(s)
}

// Fields returns a set of structured fields for logging.
func (s *Spec) Fields() log.Fields {
	return debug.StructFields(s)
}

// Status returns the current status of a spec.
// TODO: this requires a lot more thought and effort, status is complex.
func (s *Spec) Status() (string, error) {
	status := "inactive"
	// TODO: this is super naive and should be handled differently
	if services := GetServices().Filter(s.Name); len(services.Items) != 0 {
		var statuses []string
		for _, service := range services.Items {
			status, err := service.Status()
			if err != nil {
				return "", err
			}

			statuses = append(statuses, status.Status)
		}

		status = strings.Join(statuses, ", ")
	}

	return status, nil
}

// IsValid returns an error if the spec is not valid.
func (s *Spec) IsValid() error {

	// Cannot sync files, must do directories.
	fstat, err := os.Stat(s.LocalPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	if fstat != nil && !fstat.IsDir() {
		return fmt.Errorf("local path cannot be a single file, please use a directory")
	}

	// User is of format int:int
	if result, _ := regexp.MatchString("[0-9]+:[0-9]+", s.User); !result {
		return fmt.Errorf("user must be of format uid:gid, got: %s", s.User)
	}

	return nil
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
func (s *SpecMap) Create(name string, spec *Spec, force bool) error {
	if !force {
		if s.Has(name) {
			// TODO: make this into a type?
			return fmt.Errorf("name already exists")
		}

		if s.HasLike(spec) {
			return fmt.Errorf("similar spec exists")
		}
	}

	s.Items[name] = spec
	return nil
}

// Delete removes a given spec from a SpecMap
func (s *SpecMap) Delete(name string) error {
	if !s.Has(name) {
		return fmt.Errorf("does not exist")
	}

	delete(s.Items, name)
	return nil
}

// Save serializes the current SpecMap's items to the config file.
// TODO: tests:
//   missing config file
//   shorter config file (removing an entry)
func (s *SpecMap) Save() error {
	cfgPath := viper.ConfigFileUsed()
	if cfgPath == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		cfgPath = filepath.Join(home, fmt.Sprintf(".%s.yaml", "ksync"))
	}

	log.WithFields(log.Fields{
		"path": cfgPath,
	}).Debug("writing config file")

	viper.Set("spec", s.Items)
	buf, err := yaml.Marshal(viper.AllSettings())
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cfgPath, buf, 0644)
}

// HasLike checks a given spec for deep equivalence against another spec
// TODO: is this the best way to do this?
func (s *SpecMap) HasLike(target *Spec) bool {
	for _, spec := range s.Items {
		if reflect.DeepEqual(target, spec) {
			return true
		}
	}
	return false
}

// Has checks a given spec for simple equivalence against another spec
func (s *SpecMap) Has(target string) bool {
	if _, ok := s.Items[target]; ok {
		return true
	}
	return false
}
