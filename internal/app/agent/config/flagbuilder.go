package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"
)

type flags struct {
	hashKey        *string
	address        string
	cryptoKeyPath  string
	configFileName string
	pollInterval   int
	reportInterval int
	rateLimit      int
}

func (c *Config) parseFlagsVariables() *flags {
	adress := ""
	flag.StringVar(
		&adress,
		"a",
		adress,
		"address and port to run server",
	)

	cryptoKey := ""
	flag.StringVar(
		&cryptoKey,
		"crypto-key",
		cryptoKey,
		"path to the file with the public key",
	)

	key := ""
	flag.StringVar(&key, "k", "", "hash key")

	configFileName := ""
	flag.StringVar(&configFileName, "config", configFileName, "hash key")

	var durPoll int
	flag.IntVar(&durPoll, "p", defaultPollInt, "the frequency of polling metrics from the runtime package")

	var durRep int
	flag.IntVar(&durRep, "r", defaultReportInt, "the frequency of sending metrics to the server")

	var rateLimit int
	flag.IntVar(&rateLimit, "l", defaultRateLimit, "rate limit")

	flag.Parse()

	res := &flags{
		address:        adress,
		cryptoKeyPath:  cryptoKey,
		hashKey:        &key,
		pollInterval:   durPoll,
		reportInterval: durRep,
		rateLimit:      rateLimit,
		configFileName: configFileName,
	}

	return res
}

func (c *Config) setFlagsVariables(f *flags) error {
	if isFlagPassed("a") {
		c.Address = f.address
	}

	if isFlagPassed("k") {
		c.HashKey = f.hashKey
	}

	if isFlagPassed("p") {
		dur, err := time.ParseDuration(strconv.Itoa(f.pollInterval) + "s")
		if err != nil {
			return errors.New("can not parse poll_interval as duration " + err.Error())
		}
		c.PollInterval = dur
	}

	if isFlagPassed("r") {
		dur, err := time.ParseDuration(strconv.Itoa(f.reportInterval) + "s")
		if err != nil {
			return errors.New("can not parse report_interval as duration " + err.Error())
		}
		c.ReportInterval = dur
	}

	if isFlagPassed("l") {
		c.RateLimit = f.rateLimit
	}

	if isFlagPassed("crypto-key") {
		c.CryptoKeyPath = f.cryptoKeyPath
	}

	return nil
}

func (c *Config) buildFlags(f *flags) error {
	err := c.setFlagsVariables(f)
	if err != nil {
		return fmt.Errorf("failed build flags config: %w", err)
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
