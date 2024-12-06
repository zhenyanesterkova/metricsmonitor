package config

import (
	"fmt"
	"time"
)

type Config struct {
	DBConfig DataBaseConfig
	SConfig  ServerConfig
	LConfig  LoggerConfig
}

func New() *Config {
	return &Config{
		SConfig: ServerConfig{
			Address: DefaultServerAddress,
		},
		LConfig: LoggerConfig{
			Level: DefaultLogLevel,
		},
		DBConfig: DataBaseConfig{
			FileStorageConfig: &FileStorageConfig{
				FileStoragePath: DefaultFileStoragePath,
				StoreInterval:   DefaultStoreInterval * time.Second,
				Restore:         DefaultRestore,
			},
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
