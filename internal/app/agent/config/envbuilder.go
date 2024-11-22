package config

import (
	"errors"
	"os"
	"time"
)

func (c *Config) readEnv() {
	if addr, ok := os.LookupEnv("ADDRESS"); ok {
		c.addressEnv = &addr
	}
	if pollInt, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		c.pollIntervalEnv = &pollInt
	}
	if reportInt, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		c.reportIntervalEnv = &reportInt
	}
}

func (c *Config) setEnvAddress() {
	if c.addressEnv != nil {
		c.Address = *c.addressEnv
	}
}

func (c *Config) setEnvPollInterval() error {
	if c.pollIntervalEnv != nil {
		dur, err := time.ParseDuration(*c.pollIntervalEnv + "s")
		if err != nil {
			return errors.New("can not parse poll_interval as duration" + err.Error())
		}
		c.PollInterval = dur
	}
	return nil
}
func (c *Config) setEnvReportInterval() error {
	if c.reportIntervalEnv != nil {
		dur, err := time.ParseDuration(*c.reportIntervalEnv + "s")
		if err != nil {
			return errors.New("can not parse report_interval as duration" + err.Error())
		}
		c.ReportInterval = dur
	}
	return nil
}

func (c *Config) buildEnv() error {
	c.readEnv()

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
