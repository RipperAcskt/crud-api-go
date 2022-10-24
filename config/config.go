package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Addr     string
	Postgres *Postgres `yaml:"Postgres"`
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string `yaml:"DBName"`
}

func New() (*Config, error) {
	yamlFile, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("read file failed: %v", err)
	}

	c := &Config{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %v", err)
	}
	return c, nil
}

func (p *Postgres) GetConnectionUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", p.User, p.Password, p.Host, p.Port, p.DBName)
}
