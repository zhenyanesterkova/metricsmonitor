package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func (c *Config) setEnvAddress() {
	if addr, ok := os.LookupEnv("ADDRESS"); ok {
		log.Printf("env:ADDRESS=%s", addr)
		c.Address = addr
	}
}

func (c *Config) setEnvHashKey() {
	if key, ok := os.LookupEnv("KEY"); ok {
		c.HashKey = &key
	}
}

func (c *Config) setEnvPollInterval() error {
	if pollInt, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		log.Printf("env:POLL_INTERVAL=%s", pollInt)
		dur, err := time.ParseDuration(pollInt + "s")
		if err != nil {
			return errors.New("can not parse poll_interval as duration" + err.Error())
		}
		c.PollInterval = dur
	}
	return nil
}

func (c *Config) setEnvReportInterval() error {
	if reportInt, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		log.Printf("env:REPORT_INTERVAL=%s", reportInt)
		dur, err := time.ParseDuration(reportInt + "s")
		if err != nil {
			return errors.New("can not parse report_interval as duration" + err.Error())
		}
		c.ReportInterval = dur
	}
	return nil
}

func (c *Config) setEnvRateLimit() error {
	if lim, ok := os.LookupEnv("RATE_LIMIT"); ok {
		limit, err := strconv.Atoi(lim)
		if err != nil {
			return fmt.Errorf("failed parse RATE_LIMIT as int: %w", err)
		}
		c.RateLimit = limit
	}
	return nil
}

func (c *Config) buildEnv() error {
	c.setEnvAddress()

	c.setEnvHashKey()

	err := c.setEnvPollInterval()
	if err != nil {
		return err
	}

	err = c.setEnvReportInterval()
	if err != nil {
		return err
	}

	err = c.setEnvRateLimit()
	if err != nil {
		return err
	}

	return nil
}
