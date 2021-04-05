package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Platforms struct {
		Juejin struct {
			Cookie string `yaml:"cookie"`
		} `yaml:"juejin"`
		OSChina struct {
			Cookie string `yaml:"cookie"`
		} `yaml:"oschina"`
	} `yaml:"platforms"`
}

func ParseConfig(cfgFile string) (cfg *Config, err error) {
	var b []byte
	b, err = ioutil.ReadFile(cfgFile)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(b, &cfg)
	return
}
