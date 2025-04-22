package config

import (
	"fmt"
	"time"
)

type Config struct {
	DBConfig    DataBaseConfig `json:"db_config"`
	LConfig     LoggerConfig   `json:"log_config"`
	SConfig     ServerConfig   `json:"server_config"`
	RetryConfig RetryConfig    `json:"retry_config"`
}

func New() *Config {
	return &Config{
		SConfig: ServerConfig{
			Address:              DefaultServerAddress,
			CryptoPrivateKeyPath: DefaultCryptoPrivateKeyPath,
			CryptoPublicKeyPath:  DefaultCryptoPublicKeyPath,
			NeedGenKeys:          DefualtNeedGenKeys,
			ConfigsFileName:      DefaultConfigsFileName,
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
			PostgresConfig: &PostgresConfig{},
		},
		RetryConfig: RetryConfig{
			MinDelay:   DefaultMinRetryDelay,
			MaxDelay:   DefaultMaxRetryDelay,
			MaxAttempt: DefaultMaxRetryAttempt,
		},
	}
}

func (c *Config) Build() error {
	flagsVar := c.parseFlagsVariables()

	if flagsVar.config != "" {
		c.SConfig.ConfigsFileName = flagsVar.config
	}
	err := c.fileBuild()
	if err != nil {
		return fmt.Errorf("error build config from file: %w", err)
	}

	err = c.flagBuild(flagsVar)
	if err != nil {
		return fmt.Errorf("error build config from flags: %w", err)
	}

	err = c.envBuild()
	if err != nil {
		return fmt.Errorf("error build config from env: %w", err)
	}

	return nil
}
