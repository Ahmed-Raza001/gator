package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

// Config represents the structure of our config file.
type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// getConfigFilePath returns the path to the config file.
func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + "/" + configFileName, nil
}

// Read reads the config file and returns a Config struct.
func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// write writes the Config struct to the config file.
func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// SetUser updates the current user name and saves the config.
func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName

	return write(*cfg)
}
