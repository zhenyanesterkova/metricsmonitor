package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (c *Config) fileBuild() error {
	if nameConfig, ok := os.LookupEnv("CONFIG"); ok {
		c.ConfigFileName = nameConfig
	}

	currentAppPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed get path: %w", err)
	}
	configPath := filepath.Join(currentAppPath, c.ConfigFileName)
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed read config from file: %w", err)
	}

	err = json.Unmarshal(configData, c)
	if err != nil {
		return fmt.Errorf("failed unmarshal json to config: %w", err)
	}

	return nil
}
