package config

import (
	"errors"
	"os"
	"time"
)

func (c *Config) ReadEnv() {
	addr := os.Getenv("ADDRESS")
	pollInt := os.Getenv("POLL_INTERVAL")
	reportInt := os.Getenv("REPORT_INTERVAL")
	if addr != "" {
		c.addressEnv = &addr
	}
	if pollInt != "" {
		c.pollIntervalEnv = &pollInt
	}
	if reportInt != "" {
		c.reportIntervalEnv = &reportInt
	}
}

func (c *Config) SetEnvAddress() {
	if c.addressEnv != nil {
		c.Address = *c.addressEnv
	}
}

func (c *Config) SetEnvPollInterval() error {
	if c.pollIntervalEnv != nil {
		dur, err := time.ParseDuration(*c.pollIntervalEnv + "s")
		if err != nil {
			return errors.New("can not parse poll_interval as duration" + err.Error())
		}
		c.PollInterval = dur
	}
	return nil
}
func (c *Config) SetEnvReportInterval() error {
	if c.reportIntervalEnv != nil {
		dur, err := time.ParseDuration(*c.reportIntervalEnv + "s")
		if err != nil {
			return errors.New("can not parse report_interval as duration" + err.Error())
		}
		c.ReportInterval = dur
	}
	return nil
}

func (c *Config) BuildEnv() error {
	c.ReadEnv()

	c.SetEnvAddress()

	err := c.SetEnvPollInterval()
	if err != nil {
		return err
	}

	err = c.SetEnvReportInterval()
	if err != nil {
		return err
	}

	return nil
}
