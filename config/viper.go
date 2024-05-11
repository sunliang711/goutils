package config

import (
	"strings"

	"github.com/spf13/viper"
)

// newViper create an instance of viper.Viper from file [./config.(yaml|json|toml)] and env var
// envPrefix: setup env when not empty
func newViper(name, path, envPrefix string) (*viper.Viper, error) {

	vp := viper.New()

	vp.SetConfigName(name)
	vp.AddConfigPath(path)

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

type Option func(*options)

type options struct {
	configName string
	configPath string
	envPrefix  string
}

func WithConfigName(name string) Option {
	return func(o *options) {
		o.configName = name
	}
}

func WithConfigPath(path string) Option {
	return func(o *options) {
		o.configPath = path
	}
}

func WithEnvPrefix(prefix string) Option {
	return func(o *options) {
		o.envPrefix = prefix
	}
}

type Config struct {
	options
	vp *viper.Viper
}

func New(option ...Option) (*Config, error) {
	os := options{
		configName: "config",
		configPath: ".",
		envPrefix:  "APP",
	}

	for _, o := range option {
		o(&os)
	}

	vp, err := newViper(os.configName, os.configPath, os.envPrefix)
	if err != nil {
		return nil, err
	}

	return &Config{
		options: os,
		vp:      vp,
	}, nil

}

func (c *Config) Get(key string) any {
	return c.vp.Get(key)
}

func (c *Config) GetBool(key string) bool {
	return c.vp.GetBool(key)
}

func (c *Config) GetInt(key string) int {
	return c.vp.GetInt(key)
}

func (c *Config) GetFloat64(key string) float64 {
	return c.vp.GetFloat64(key)
}

func (c *Config) GetString(key string) string {
	return c.vp.GetString(key)
}
