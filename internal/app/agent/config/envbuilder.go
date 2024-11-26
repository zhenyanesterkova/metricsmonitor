package config

import (
	"errors"
	"os"
	"time"
)

func (c *Config) setEnvAddress() {
	if addr, ok := os.LookupEnv("ADDRESS"); ok {
		c.Address = addr
	}
}

func (c *Config) setEnvPollInterval() error {
	if pollInt, ok := os.LookupEnv("POLL_INTERVAL"); ok {
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
		dur, err := time.ParseDuration(reportInt + "s")
		if err != nil {
			return errors.New("can not parse report_interval as duration" + err.Error())
		}
		c.ReportInterval = dur
	}
	return nil
}

func (c *Config) buildEnv() error {
	c.setEnvAddress()

	err := c.setEnvPollInterval()
	if err != nil {
		return err
	}

	err = c.setEnvReportInterval()
	if err != nil {
		return err
	}

	return nil
}
