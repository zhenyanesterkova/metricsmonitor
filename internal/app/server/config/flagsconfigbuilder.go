package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"
)

type flags struct {
	hashKey         *string
	fileStoragePath string
	dsn             string
	tempDur         int
	restore         bool
}

func (c *Config) parseFlagsVariables() flags {
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

	flag.StringVar(
		&c.SConfig.CryptoPrivateKeyPath,
		"crypto-key",
		c.SConfig.CryptoPrivateKeyPath,
		"path to the file with the private key",
	)

	flag.StringVar(
		&c.SConfig.CryptoPublicKeyPath,
		"crypto-pub-key",
		c.SConfig.CryptoPublicKeyPath,
		"path to the file with the private key",
	)

	flag.Parse()

	res := flags{
		fileStoragePath: fileStoragePath,
		tempDur:         tempDur,
		restore:         restore,
		dsn:             dsn,
		hashKey:         &hashKey,
	}
	return res
}

func (c *Config) setFlagsVariables(f flags) error {
	if isFlagPassed("f") {
		c.DBConfig.FileStorageConfig.FileStoragePath = f.fileStoragePath
	}

	if isFlagPassed("i") {
		dur, err := time.ParseDuration(strconv.Itoa(f.tempDur) + "s")
		if err != nil {
			return errors.New("can not parse store interval as duration " + err.Error())
		}
		c.DBConfig.FileStorageConfig.StoreInterval = dur
	}

	if isFlagPassed("r") {
		c.DBConfig.FileStorageConfig.Restore = f.restore
	}

	if isFlagPassed("d") {
		if c.DBConfig.PostgresConfig == nil {
			c.DBConfig.PostgresConfig = &PostgresConfig{}
		}
		c.DBConfig.PostgresConfig.DSN = f.dsn
	}

	if isFlagPassed("k") {
		c.SConfig.HashKey = f.hashKey
	}

	return nil
}

func (c *Config) flagBuild() error {
	flagsVar := c.parseFlagsVariables()
	err := c.setFlagsVariables(flagsVar)
	if err != nil {
		return fmt.Errorf("failed set cfg from flags var: %w", err)
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
