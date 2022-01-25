package config

import (
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/juju/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Platforms Platforms `yaml:"platforms,omitempty"`
}

type Platforms struct {
	Juejin  Juejin  `yaml:"juejin,omitempty"`
	OSChina OSChina `yaml:"oschina,omitempty"`
	Github  Github  `yaml:"github,omitempty"`
}

type Juejin struct {
	Cookie string `yaml:"cookie,omitempty"`
}

type OSChina struct {
	Cookie string `yaml:"cookie,omitempty"`
}

type Github struct {
	Token string `yaml:"token,omitempty"`
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

func GetConfigDir() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", errors.Trace(err)
	}

	cfgDir := filepath.Join(homeDir, ".config", "articli")
	if err = os.MkdirAll(cfgDir, os.ModePerm); err != nil {
		return "", errors.Trace(err)
	}
	return cfgDir, nil
}
