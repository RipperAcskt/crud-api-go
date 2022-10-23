package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Url  string
	Ip   string
	Port string
}

func New() (*Config, error) {
	yamlFile, err := ioutil.ReadFile("../config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("readFile failed: %v", err)
	}

	c := &Config{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %v", err)
	}
	return c, nil
}