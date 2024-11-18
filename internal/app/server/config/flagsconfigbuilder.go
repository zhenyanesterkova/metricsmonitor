package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"
)

func (c *Config) setFlagsVariables() {
	var (
		address       string
		lvllog        string
		storeInterval int
		filePath      string
		restore       bool
	)
	c.SConfig.addressParam = &address
	c.LConfig.levelParam = &lvllog
	c.RConfig.storeIntervalParamFl = &storeInterval
	c.RConfig.filePathParam = &filePath
	c.RConfig.restoreParam = &restore

	flag.StringVar(c.SConfig.addressParam, "a", DefaultServerAddress, "address and port to run server")
	flag.StringVar(c.LConfig.levelParam, "l", DefaultLogLevel, "log level")
	flag.IntVar(c.RConfig.storeIntervalParamFl, "i", DefaultStoreInterval, "store interval")
	flag.StringVar(c.RConfig.filePathParam, "f", DefaultFileStoragePath, "file storage path")
	flag.BoolVar(c.RConfig.restoreParam, "r", DefaultRestore, "need restore")
	flag.Parse()
}

func (c *Config) setFlServerConfig() {
	if c.SConfig.addressParam != nil {
		c.SConfig.Address = *c.SConfig.addressParam
	}
}

func (c *Config) setFlLoggerConfig() {
	if c.LConfig.levelParam != nil {
		c.LConfig.Level = *c.LConfig.levelParam
	}
}

func (c *Config) setFlRestoreConfig() error {
	if c.RConfig.storeIntervalParamFl != nil {
		dur, err := time.ParseDuration(strconv.Itoa(*c.RConfig.storeIntervalParamFl) + "s")
		if err != nil {
			return errors.New("can not parse restore interval as duration " + err.Error())
		}
		c.RConfig.StoreInterval = dur
	}

	if c.RConfig.filePathParam != nil {
		c.RConfig.FileStoragePath = *c.RConfig.filePathParam
	}

	if c.RConfig.restoreParam != nil {
		c.RConfig.Restore = *c.RConfig.restoreParam
	}

	return nil
}

func (c *Config) flagBuild() error {
	c.setFlagsVariables()
	c.setFlServerConfig()
	c.setFlLoggerConfig()
	err := c.setFlRestoreConfig()
	if err != nil {
		return fmt.Errorf("build flag config error: %w", err)
	}
	return nil
}
