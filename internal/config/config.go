package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string
}

const configFileName = ".gatorconfig.json"

// set path to config file
func GetConfigFilePath() (string, error) {

	configPath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(configPath, configFileName)

	return fullPath, nil
}

// read config file
func Read() (Config, error) {
	cfgPath, err := GetConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	cfgBytes, err := os.ReadFile(cfgPath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// set current user
func (c *Config) SetUser(user string) error {

	c.CurrentUserName = user
	return write(*c)
}

// update config file with current user
func write(cfg Config) error {

	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	fileName, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
