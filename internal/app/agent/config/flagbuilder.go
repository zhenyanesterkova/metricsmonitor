package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"
)

func (c *Config) setFlags() error {
	flag.StringVar(&c.Address, "a", c.Address, "address and port to run server")
	flag.StringVar(&c.CryptoKeyPath, "crypto-key", c.CryptoKeyPath, "path to the file with the public key")

	var key string
	flag.StringVar(&key, "k", "", "hash key")

	var durPoll int
	flag.IntVar(&durPoll, "p", defaultPollInt, "the frequency of polling metrics from the runtime package")

	var durRep int
	flag.IntVar(&durRep, "r", defaultReportInt, "the frequency of sending metrics to the server")

	var rateLimit int
	flag.IntVar(&rateLimit, "l", defaultRateLimit, "rate limit")

	flag.Parse()

	if isFlagPassed("k") {
		c.HashKey = &key
	}

	if isFlagPassed("p") {
		dur, err := time.ParseDuration(strconv.Itoa(durPoll) + "s")
		if err != nil {
			return errors.New("can not parse poll_interval as duration " + err.Error())
		}
		c.PollInterval = dur
	}

	if isFlagPassed("r") {
		dur, err := time.ParseDuration(strconv.Itoa(durRep) + "s")
		if err != nil {
			return errors.New("can not parse report_interval as duration " + err.Error())
		}
		c.ReportInterval = dur
	}

	if isFlagPassed("l") {
		c.RateLimit = rateLimit
	}

	return nil
}

func (c *Config) buildFlags() error {
	err := c.setFlags()
	if err != nil {
		return fmt.Errorf("config func buildFlags(): %w", err)
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
