package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func (c *Config) setEnvServerConfig() {
	if envEndpoint, ok := os.LookupEnv("ADDRESS"); ok {
		log.Printf("env:ADDRESS=%s", envEndpoint)
		c.SConfig.Address = envEndpoint
	}
}

func (c *Config) setEnvLoggerConfig() {
	if envLogLevel, ok := os.LookupEnv("LOG_LEVEL"); ok {
		log.Printf("env:LOG_LEVEL=%s", envLogLevel)
		c.LConfig.Level = envLogLevel
	}
}

func (c *Config) setEnvRestoreConfig() error {
	if envStoreInt, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		log.Printf("env:STORE_INTERVAL=%s", envStoreInt)
		dur, err := time.ParseDuration(envStoreInt + "s")
		if err != nil {
			return errors.New("can not parse store interval as duration" + err.Error())
		}
		c.RConfig.StoreInterval = dur
	}

	if envFileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		log.Printf("env:FILE_STORAGE_PATH=%s", envFileStoragePath)
		c.RConfig.FileStoragePath = envFileStoragePath
	}

	if envRestore, ok := os.LookupEnv("RESTORE"); ok {
		log.Printf("env:RESTORE=%s", envRestore)
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
