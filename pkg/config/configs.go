package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	DownloadFolder     string     `yaml:"download-folder"`
	TmpFolder          string     `yaml:"temp-folder"`
	ConcurrentDownload uint       `yaml:"concurrent-download"`
	Items              []ItemInfo `yaml:"items"`
}

type ItemInfo struct {
	Name  string                 `yaml:"name"`
	Type  string                 `yaml:"type"`
	Info  string                 `yaml:"info"`
	Extra map[string]interface{} `yaml:"extra"`
}

func ReadConfig(configPath string) (*Config, error) {
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
