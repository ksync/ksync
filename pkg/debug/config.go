package debug

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/go-yaml/yaml"
	log "github.com/sirupsen/logrus"
)

var (
	mismatchError = "Detected mismatched type in config. This may be to due to updates or corruption in your config file. Got %v but wanted %v"
)

type ConfigFormat struct {
	APIKey     string `yaml:"apikey"`
	Context    string `yaml:"context"`
	DockerRoot string `yaml:"docker-root"`
	LogLevel   string `yaml:"log-level"`
	Namespace  string `yaml:"namespace"`
	Output     string `yaml:"output"`
	Port       int    `yaml:"port"`
	Spec struct {
		Name           string `yaml:"name"`
		ContainerName  string              `yaml:"containername"`
		Pod            string              `yaml:"pod"`
		Selector       []map[string][]string `yaml:"selector"`
		Namespace      string              `yaml:"namespace"`
		LocalPath      string              `yaml:"localpath"`
		RemotePath     string              `yaml:"remotepath"`
		Reload         string              `yaml:"reload"`
		LocalReadOnly  string              `yaml:"localreadonly"`
		RemoteReadOnly string              `yaml:"remotereadonly"`
		SyncthingPort string `yaml:"syncthing-port"`
	}
}

func ValidateConfig(path string) error {
	config, err := ReadInConfig(path)
	if err != nil {
		return err
	}

	err = CheckOldTypes(config)
	if err != nil {
		return err
	}

	return nil
}

func ReadInConfig(path string) (*ConfigFormat, error) {
	config := ConfigFormat{}
	file, err := ioutil.ReadFile(fmt.Sprintf("%v/%s", path, "ksync.yaml"))
	log.Info(fmt.Sprintf("%v/%s", path, "ksync.yaml"))
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func CheckOldTypes(config *ConfigFormat) error {
	expectedConfig := ConfigFormat{}
	log.Infof("real: %+v, expected: %+v", config, expectedConfig)

	// Check `selectors` is not a key-value (>0.2.6)
	if reflect.TypeOf(config.Spec.Selector) != reflect.TypeOf(expectedConfig.Spec.Selector) {
		log.Errorf(mismatchError, reflect.TypeOf(config.Spec.Selector), reflect.TypeOf(expectedConfig.Spec.Selector))
	}

	return nil
}
