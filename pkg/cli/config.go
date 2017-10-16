package cli

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// InitConfig constructs the configuration from a local configuration file
// or environment variables if available. This is placed in the global `viper`
// instance.
func InitConfig(name string) {
	viper.SupportedExts = []string{"yaml", "yml"}

	if viper.GetString("config") != "" {
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(fmt.Sprintf(".%s", name))

		cfgPath := filepath.Join(home, fmt.Sprintf(".%s.yaml", name))
		fobj, _ := os.OpenFile(cfgPath, os.O_CREATE|os.O_WRONLY, 0644)
		fobj.Close()
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// TODO: the level here is *always* the default, need a better solution
		// for output.
		log.WithFields(log.Fields{
			"file": viper.ConfigFileUsed(),
		}).Debug("using config file")
	}
}
