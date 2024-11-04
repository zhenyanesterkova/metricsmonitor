package config

import (
	"time"
)

type Config struct {
	Address           string
	PollInterval      time.Duration
	ReportInterval    time.Duration
	addressFl         *string
	pollIntervalFl    *int
	reportIntervalFl  *int
	addressEnv        *string
	pollIntervalEnv   *string
	reportIntervalEnv *string
}

func New() *Config {
	return &Config{
		Address:        "localhost:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
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
