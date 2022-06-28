package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Authorization struct {
	Method    string   `yaml:"method"`
	PathRegEx []string `yaml:"pathregex"`
}

type requesAuthz struct {
	Method    string   `json:method`
	PathRegEx []string `json:pathregex`
}

type headerRules map[string][]requesAuthz

//type ConfigFile map[string][]Authorization
type ConfigFile struct {
	Listen         string      `yaml:"listen"`
	Pdns_api_url   string      `yaml:"pdns_api_url"`
	Pdns_api_token string      `yaml:"pdns_api_token"`
	DebugMode      bool        `yaml:"debugMode"`
	Rules          headerRules `yaml:"rules"`
}

var (
	config ConfigFile
)

func (cfg *ConfigFile) loadConfig(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		return err
	}
	return nil
}
