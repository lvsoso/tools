package config

import (
	"errors"
	"os"
	"time"

	"github.com/spf13/viper"
)

var (
	Config            ConfigConf
	ErrConfigNotFound = errors.New("config not found")
)

type ConfigConf struct {
	GitUsername string         `yaml:"gitUsername,omitempty"`
	GitPassword string         `yaml:"gitPassword,omitempty"`
	Targets     []TargetConfig `yaml:"targets,omitempty"`
	mapTargets  map[string]TargetConfig
}

func (cc ConfigConf) AddConfig(name string, config TargetConfig) {
	cc.mapTargets[name] = config
}

func (cc ConfigConf) GetConfig(name string) (*TargetConfig, error) {
	if cfg, ok := cc.mapTargets[name]; ok {
		return &cfg, nil
	}
	return nil, ErrConfigNotFound
}

type TargetConfig struct {
	Name          string            `yaml:"name,omitempty"`
	Driver        string            `yaml:"driver,omitempty"`
	Root          string            `yaml:"root,omitempty"`
	URL           string            `yaml:"url,omitempty"`
	BackendConfig map[string]string `yaml:"backendConfig,omitempty"`
	Timeout       time.Duration     `yaml:"timeout,omitempty"`
}

func (tc TargetConfig) Get(key string) (value string, ok bool) {
	v, ok := tc.BackendConfig[key]
	return v, ok
}

func init() {
	configPath := os.Getenv("QUQ_CONFIG_PATH")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	Config = ConfigConf{
		mapTargets: map[string]TargetConfig{},
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}

	for _, tc := range Config.Targets {
		Config.mapTargets[tc.Name] = tc
	}
}
