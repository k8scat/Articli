package config

import (
	"io/ioutil"

	"github.com/juju/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Platforms Platforms `yaml:"platforms,omitempty"`
}

type Platforms struct {
	Juejin  Juejin  `yaml:"juejin,omitempty"`
	OSChina OSChina `yaml:"oschina,omitempty"`
}

type Juejin struct {
	Cookie string `yaml:"cookie,omitempty"`
}

type OSChina struct {
	Cookie string `yaml:"cookie,omitempty"`
}

func ParseConfig(cfgFile string) (*Config, error) {
	b, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, errors.Trace(err)
	}
	cfg := new(Config)
	err = yaml.Unmarshal(b, &cfg)
	return cfg, errors.Trace(err)
}

func SaveConfig(cfgFile string, cfg *Config) error {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Trace(err)
	}
	err = ioutil.WriteFile(cfgFile, b, 0644)
	return errors.Trace(err)
}
