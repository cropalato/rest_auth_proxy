package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

type Authorization struct {
	Method    string   `yaml:"method"`
	PathRegEx []string `yaml:"pathregex"`
}

type requesAuthz struct {
	Method    string   `json:"method"`
	PathRegEx []string `json:"pathregex"`
}

type headerRules map[string][]requesAuthz

//type ConfigFile map[string][]Authorization
type ConfigFile struct {
	Listen           string      `yaml:"listen"`
	Server_api_url   string      `yaml:"server_api_url"`
	Server_api_token string      `yaml:"server_api_token"`
	Header_token     string      `yaml:"header_token"`
	Rules            headerRules `yaml:"rules"`
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
