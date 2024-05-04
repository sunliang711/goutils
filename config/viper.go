package config

import (
	"strings"

	"github.com/spf13/viper"
)

var Viper *viper.Viper

func init() {
	Viper, _ = NewViper("APP")
}

// NewViper create an instance of viper.Viper from file [./config.(yaml|json|toml)] and env var
// envPrefix: setup env when not empty
func NewViper(envPrefix string) (*viper.Viper, error) {

	vp := viper.New()

	vp.SetConfigName("config")
	vp.AddConfigPath(".")

	if envPrefix != "" {
		vp.SetEnvPrefix(envPrefix)
		vp.AutomaticEnv()
		vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	}

	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return vp, nil
}
