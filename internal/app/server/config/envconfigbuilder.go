package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

func (c *Config) setEnvServerConfig() {
	envEndpoint := os.Getenv("ADDRESS")

	if envEndpoint != "" {
		c.SConfig.Address = envEndpoint
	}
}

func (c *Config) setEnvLoggerConfig() {
	envLogLevel := os.Getenv("LOG_LEVEL")

	if envLogLevel != "" {
		c.LConfig.Level = envLogLevel
	}
}

func (c *Config) setEnvRestoreConfig() error {
	envStoreInt := os.Getenv("STORE_INTERVAL")

	if envStoreInt != "" {
		dur, err := time.ParseDuration(envStoreInt + "s")
		if err != nil {
			return errors.New("can not parse store interval as duration" + err.Error())
		}
		c.RConfig.StoreInterval = dur
	}

	envFileStoragePath := os.Getenv("FILE_STORAGE_PATH")

	if envFileStoragePath != "" {
		c.RConfig.FileStoragePath = envFileStoragePath
	}

	envRestore := os.Getenv("RESTORE")

	if envRestore != "" {
		path, err := strconv.ParseBool(envRestore)
		if err != nil {
			return errors.New("can not parse need store" + err.Error())
		}
		c.RConfig.Restore = path
	}

	return nil
}

func (c *Config) envBuild() error {
	c.setEnvServerConfig()
	c.setEnvLoggerConfig()
	err := c.setEnvRestoreConfig()
	if err != nil {
		return fmt.Errorf("build env config error: %w", err)
	}
	return nil
}
