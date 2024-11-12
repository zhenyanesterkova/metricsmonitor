package config

import (
	"time"
)

const (
	defaultAddress   = "localhost:8080"
	defaultPollInt   = 2
	defaultReportInt = 10
)

type Config struct {
	pollIntervalFl    *int
	reportIntervalFl  *int
	addressEnv        *string
	pollIntervalEnv   *string
	reportIntervalEnv *string
	addressFl         *string
	Address           string
	PollInterval      time.Duration
	ReportInterval    time.Duration
}

func New() *Config {
	return &Config{
		Address:        defaultAddress,
		PollInterval:   defaultPollInt * time.Second,
		ReportInterval: defaultReportInt * time.Second,
	}
}

func (c *Config) Build() error {
	err := c.BuildFlags()
	if err != nil {
		return err
	}
	err = c.BuildEnv()
	if err != nil {
		return err
	}

	return nil
}
