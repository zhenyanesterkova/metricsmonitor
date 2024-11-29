package config

import (
	"fmt"
	"time"
)

type Config struct {
	SConfig ServerConfig
	LConfig LoggerConfig
	RConfig RestoreConfig
}

func New() *Config {
	return &Config{
		SConfig: ServerConfig{
			Address: DefaultServerAddress,
		},
		LConfig: LoggerConfig{
			Level: DefaultLogLevel,
		},
		RConfig: RestoreConfig{
			FileStoragePath: DefaultFileStoragePath,
			Restore:         DefaultRestore,
			StoreInterval:   DefaultStoreInterval * time.Second,
		},
	}
}

func (c *Config) Build() error {
	err := c.flagBuild()
	if err != nil {
		return fmt.Errorf("error build config from flags: %w", err)
	}

	err = c.envBuild()
	if err != nil {
		return fmt.Errorf("error build config from env: %w", err)
	}

	return nil
}
