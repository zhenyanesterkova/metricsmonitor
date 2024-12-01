package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"
)

func (c *Config) setFlagsVariables() error {
	flag.StringVar(&c.SConfig.Address, "a", c.SConfig.Address, "address and port to run server")
	flag.StringVar(&c.LConfig.Level, "l", c.LConfig.Level, "log level")

	var tempDur int
	flag.IntVar(&tempDur, "i", DefaultStoreInterval, "store interval")

	flag.StringVar(&c.DBConfig.FileStoragePath, "f", c.DBConfig.FileStoragePath, "file storage path")
	flag.BoolVar(&c.DBConfig.Restore, "r", c.DBConfig.Restore, "need restore")
	flag.StringVar(&c.DBConfig.DSN, "d", c.DBConfig.DSN, "database dsn")
	flag.Parse()

	if isFlagPassed("f") || isFlagPassed("r") {
		c.DBConfig.DBType = FileStorageType
	}

	if isFlagPassed("d") {
		c.DBConfig.DBType = PostgresStorageType
	}

	if isFlagPassed("i") {
		dur, err := time.ParseDuration(strconv.Itoa(tempDur) + "s")
		if err != nil {
			return errors.New("can not parse store interval as duration " + err.Error())
		}
		c.DBConfig.StoreInterval = dur
	}

	return nil
}

func (c *Config) flagBuild() error {
	err := c.setFlagsVariables()
	if err != nil {
		return fmt.Errorf("config func flagBuild(): %w", err)
	}
	return nil
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
