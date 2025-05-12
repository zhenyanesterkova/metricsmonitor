package config

import (
	"fmt"
	"net"
	"time"
)

type Config struct {
	SConfig     ServerConfig   `json:"server_config"`
	DBConfig    DataBaseConfig `json:"db_config"`
	LConfig     LoggerConfig   `json:"log_config"`
	RetryConfig RetryConfig    `json:"retry_config"`
}

func New() *Config {
	return &Config{
		SConfig: ServerConfig{
			Address:              DefaultServerAddress,
			CryptoPrivateKeyPath: DefaultCryptoPrivateKeyPath,
			CryptoPublicKeyPath:  DefaultCryptoPublicKeyPath,
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

	err = c.setTrustedSubnet()
	if err != nil {
		return fmt.Errorf("error build config: %w", err)
	}

	return nil
}

func (c *Config) setTrustedSubnet() error {
	if c.SConfig.StringCIDR == "" {
		return nil
	}

	_, ipNet, err := net.ParseCIDR(c.SConfig.StringCIDR)
	if err != nil {
		return fmt.Errorf("failed parse trusted subnet from string (%s): %w", c.SConfig.StringCIDR, err)
	}
	c.SConfig.TIpNet = ipNet
	return nil
}
