package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"
)

type flags struct {
	adress          string
	config          string
	logLevel        string
	cryptoKey       string
	cryptoPublicKey string
	hashKey         *string
	fileStoragePath string
	dsn             string
	tempDur         int
	restore         bool
}

func (c *Config) parseFlagsVariables() *flags {
	adress := ""
	flag.StringVar(
		&adress,
		"a",
		adress,
		"address and port to run server",
	)

	config := ""
	flag.StringVar(
		&config,
		"c",
		config,
		"config name",
	)

	logLevel := ""
	flag.StringVar(
		&logLevel,
		"l",
		logLevel,
		"log level",
	)

	var tempDur int
	flag.IntVar(
		&tempDur,
		"i",
		tempDur,
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

	cryptoKey := ""
	flag.StringVar(
		&cryptoKey,
		"crypto-key",
		cryptoKey,
		"path to the file with the private key",
	)

	cryptoPublicKey := ""
	flag.StringVar(
		&cryptoPublicKey,
		"crypto-pub-key",
		cryptoPublicKey,
		"path to the file with the private key",
	)

	flag.Parse()

	res := &flags{
		adress:          adress,
		config:          config,
		logLevel:        logLevel,
		cryptoKey:       cryptoKey,
		cryptoPublicKey: cryptoPublicKey,
		fileStoragePath: fileStoragePath,
		tempDur:         tempDur,
		restore:         restore,
		dsn:             dsn,
		hashKey:         &hashKey,
	}
	return res
}

func (c *Config) setFlagsVariables(f *flags) error {
	if isFlagPassed("a") {
		c.SConfig.Address = f.adress
	}
	if isFlagPassed("c") {
		c.SConfig.ConfigsFileName = f.config
	}
	if isFlagPassed("l") {
		c.LConfig.Level = f.logLevel
	}
	if isFlagPassed("crypto-key") {
		c.SConfig.CryptoPrivateKeyPath = f.cryptoKey
	}
	if isFlagPassed("crypto-pub-key") {
		c.SConfig.CryptoPublicKeyPath = f.cryptoPublicKey
	}
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

func (c *Config) flagBuild(flagsVar *flags) error {
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
