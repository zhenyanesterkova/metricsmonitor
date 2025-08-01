package config

import (
	"time"
)

const (
	defaultAddress        = "localhost:8080"
	defaultGRPCAddress    = "localhost:8888"
	defaultPollInt        = 2
	defaultReportInt      = 10
	defaultRateLimit      = 3
	defaultCryptoKeyPath  = "example-public.crt"
	DefaultCertPath       = "agent.crt"
	defaultConfigFileName = "agent_config.json"
)

type Config struct {
	HashKey        *string       `json:"hash_key"`
	Address        string        `json:"address"`
	GRPCAddress    string        `json:"grpc_address"`
	CryptoKeyPath  string        `json:"crypto_key"`
	CertPath       string        `json:"cert_path"`
	ConfigFileName string        `json:"config"`
	PollInterval   time.Duration `json:"poll_interval"`
	ReportInterval time.Duration `json:"report_interval"`
	RateLimit      int           `json:"rate_limit"`
}

func New() *Config {
	return &Config{
		Address:        defaultAddress,
		GRPCAddress:    defaultGRPCAddress,
		PollInterval:   defaultPollInt * time.Second,
		ReportInterval: defaultReportInt * time.Second,
		RateLimit:      defaultRateLimit,
		CryptoKeyPath:  defaultCryptoKeyPath,
		CertPath:       DefaultCertPath,
		ConfigFileName: defaultConfigFileName,
	}
}

func (c *Config) Build() error {
	flags := c.parseFlagsVariables()

	if flags.configFileName != "" {
		c.ConfigFileName = flags.configFileName
	}
	err := c.fileBuild()
	if err != nil {
		return err
	}

	err = c.buildFlags(flags)
	if err != nil {
		return err
	}

	err = c.buildEnv()
	if err != nil {
		return err
	}

	return nil
}
