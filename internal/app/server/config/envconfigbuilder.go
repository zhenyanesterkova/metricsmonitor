package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

func (c *Config) setEnvServerConfig() {
	if envEndpoint, ok := os.LookupEnv("ADDRESS"); ok {
		c.SConfig.Address = envEndpoint
	}
	if envHashKey, ok := os.LookupEnv("KEY"); ok {
		c.SConfig.HashKey = &envHashKey
	}
}

func (c *Config) setEnvLoggerConfig() {
	if envLogLevel, ok := os.LookupEnv("LOG_LEVEL"); ok {
		c.LConfig.Level = envLogLevel
	}
}

func (c *Config) setDBConfig() error {
	if envFileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		c.DBConfig.FileStorageConfig.FileStoragePath = envFileStoragePath
	}

	if envStoreInt, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		dur, err := time.ParseDuration(envStoreInt + "s")
		if err != nil {
			return errors.New("can not parse store interval as duration" + err.Error())
		}
		c.DBConfig.FileStorageConfig.StoreInterval = dur
	}

	if envRestore, ok := os.LookupEnv("RESTORE"); ok {
		restore, err := strconv.ParseBool(envRestore)
		if err != nil {
			return errors.New("can not parse need store" + err.Error())
		}
		c.DBConfig.FileStorageConfig.Restore = restore
	}

	if dsn, ok := os.LookupEnv("DATABASE_DSN"); ok {
		if c.DBConfig.PostgresConfig == nil {
			c.DBConfig.PostgresConfig = &PostgresConfig{}
		}
		c.DBConfig.PostgresConfig.DSN = dsn
	}

	return nil
}

func (c *Config) envBuild() error {
	c.setEnvServerConfig()
	c.setEnvLoggerConfig()
	err := c.setDBConfig()
	if err != nil {
		return fmt.Errorf("build env config error: %w", err)
	}
	return nil
}
