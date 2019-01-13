package torpoller

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var configPath string

type Config struct {
	DownloadFolder string                 `yaml:"download-folder"`
	Items          []ItemInfo             `yaml:"items"`
}

type ItemInfo struct {
	Name  string                 `yaml:"name"`
	Type  string                 `yaml:"type"`
	Info  string                 `yaml:"info"`
	Extra map[string]interface{} `yaml:"extra"`
}

func SetConfigPath(path string) {
	configPath = path
}

func ReadConfig() (*Config, error) {
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file, cause: %s", err.Error())
	}
	cfg := &Config{}
	err = yaml.Unmarshal(bytes, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file, cause: %s", err.Error())
	}
	return cfg, nil
}

