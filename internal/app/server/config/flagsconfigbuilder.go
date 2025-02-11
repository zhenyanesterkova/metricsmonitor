package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"
)

func (c *Config) setFlagsVariables() error {
	flag.StringVar(
		&c.SConfig.Address,
		"a",
		c.SConfig.Address,
		"address and port to run server",
	)

	flag.StringVar(
		&c.LConfig.Level,
		"l",
		c.LConfig.Level,
		"log level",
	)

	var tempDur int
	flag.IntVar(
		&tempDur,
		"i",
		DefaultStoreInterval,
		"store interval",
	)

	fileStoragePath := ""
	flag.StringVar(
		&fileStoragePath,
		"f",
		fileStoragePath,
		"file storage path",
	)

	restore := false
	flag.BoolVar(
		&restore,
		"r",
		restore,
		"need restore",
	)

	dsn := ""
	flag.StringVar(
		&dsn,
		"d",
		dsn,
		"database dsn",
	)

	hashKey := ""
	flag.StringVar(
		&hashKey,
		"k",
		hashKey,
		"hash key",
	)

	flag.Parse()

	if isFlagPassed("f") {
		c.DBConfig.FileStorageConfig.FileStoragePath = fileStoragePath
	}

	if isFlagPassed("i") {
		dur, err := time.ParseDuration(strconv.Itoa(tempDur) + "s")
		if err != nil {
			return errors.New("can not parse store interval as duration " + err.Error())
		}
		c.DBConfig.FileStorageConfig.StoreInterval = dur
	}

	if isFlagPassed("r") {
		c.DBConfig.FileStorageConfig.Restore = restore
	}

	if isFlagPassed("d") {
		if c.DBConfig.PostgresConfig == nil {
			c.DBConfig.PostgresConfig = &PostgresConfig{}
		}
		c.DBConfig.PostgresConfig.DSN = dsn
	}

	if isFlagPassed("k") {
		c.SConfig.HashKey = &hashKey
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
