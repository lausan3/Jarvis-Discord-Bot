package environment

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configuration struct {
	Discord DiscordConfig `mapstructure:",squash"`
}

func LoadEnvVars(config *Configuration) error {
	logrus.Info("Loading environment variables...")

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("Error reading config file, %s", err)
		return err
	}

	if err := viper.Unmarshal(&config); err != nil {
		logrus.Errorf("Unable to decode into struct, %v", err)
		return err
	}

	return nil
}
