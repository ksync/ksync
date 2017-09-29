package cli

import (
	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig(name string) {
	if viper.GetString("config") != "" {
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalf("%v", err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(fmt.Sprintf(".%s", name))
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
