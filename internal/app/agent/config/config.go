package config

import (
	"time"
)

const (
	defaultAddress       = "localhost:8080"
	defaultPollInt       = 2
	defaultReportInt     = 10
	defaultRateLimit     = 3
	defaultCryptoKeyPath = "../../../../build/crypto/pub"
)

type Config struct {
	HashKey        *string
	Address        string
	CryptoKeyPath  string
	PollInterval   time.Duration
	ReportInterval time.Duration
	RateLimit      int
}

func New() *Config {
	return &Config{
		Address:        defaultAddress,
		PollInterval:   defaultPollInt * time.Second,
		ReportInterval: defaultReportInt * time.Second,
		RateLimit:      defaultRateLimit,
		CryptoKeyPath:  defaultCryptoKeyPath,
	}
}

func (c *Config) Build() error {
	err := c.buildFlags()
	if err != nil {
		return err
	}

	err = c.buildEnv()
	if err != nil {
		return err
	}

	return nil
}
