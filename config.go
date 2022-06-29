package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

// Authorization is used to validate it method + URL is allowed.
type Authorization struct {
	Method    string   `yaml:"method"`
	PathRegEx []string `yaml:"pathregex"`
}

type requesAuthz struct {
	Method    string   `json:"method"`
	PathRegEx []string `json:"pathregex"`
}

type headerRules map[string][]requesAuthz

//ConfigFile is the structure used to map the config file.
type ConfigFile struct {
	Listen         string      `yaml:"listen"`
	ServerAPIURL   string      `yaml:"serverAPIURL"`
	ServerAPIToken string      `yaml:"serverAPIToken"`
	HeaderToken    string      `yaml:"headerToken"`
	Rules          headerRules `yaml:"rules"`
}

var (
	config ConfigFile
)

func (cfg *ConfigFile) loadConfig(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		klog.Warning(fmt.Sprintf("Config file not found on %v.", path))
	}

	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		return err
	}
	return nil
}
