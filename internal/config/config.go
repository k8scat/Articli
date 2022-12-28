package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"

	"github.com/juju/errors"
	"gopkg.in/yaml.v3"
)

var (
	CfgFile string
	Cfg     *Config
)

type Config struct {
	Auth map[string]string `yaml:"auth"`
}

func (c *Config) SetAuth(name, rawAuth string) {
	if c.Auth == nil {
		c.Auth = make(map[string]string)
	}
	c.Auth[name] = rawAuth
}

func Parse() error {
	c, err := parse(CfgFile)
	if err != nil {
		return err
	}
	Cfg = c
	return nil
}

func Save() error {
	return save(CfgFile, Cfg)
}

func parse(f string) (*Config, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, errors.Trace(err)
	}
	cfg := new(Config)
	err = yaml.Unmarshal(b, &cfg)
	return cfg, errors.Trace(err)
}

func save(f string, cfg *Config) error {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Trace(err)
	}
	err = os.WriteFile(f, b, 0644)
	return errors.Trace(err)
}

func GetConfigDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Printf("warning: get home dir failed: %s\n", err)
		return ""
	}

	cfgDir := filepath.Join(homeDir, ".config", "articli")
	if err = os.MkdirAll(cfgDir, os.ModePerm); err != nil {
		fmt.Printf("warning: create config dir failed: %s\n", err)
		return ""
	}
	return cfgDir
}
