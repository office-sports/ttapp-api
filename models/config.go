package models

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// DBConfig struct config
type DBConfig struct {
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
}

// APIConfig struct config
type APIConfig struct {
	Port int `yaml:"port"`
}

type MessageConfig struct {
	Hook      string `yaml:"hook"`
	ChannelId string `yaml:"channel_id"`
}

type Frontend struct {
	Url string `yaml:"url"`
}

// Config type
type Config struct {
	DBConfig      DBConfig      `yaml:"db"`
	APIConfig     APIConfig     `yaml:"api"`
	MessageConfig MessageConfig `yaml:"message"`
	Frontend      Frontend      `yaml:"frontend"`
}

// GetConfig loads configuration from yaml file
func GetConfig(configFileParam string) (*Config, error) {
	// Default location
	configFilePath := "config/config.yaml"
	if len(configFileParam) > 0 {
		configFilePath = configFileParam
	}
	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetMysqlConnectionString returns mysql connection string
func (config *Config) GetMysqlConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@(%s:%d)/%s",
		config.DBConfig.Username,
		config.DBConfig.Password,
		config.DBConfig.Address,
		config.DBConfig.Port,
		config.DBConfig.Schema)
}
