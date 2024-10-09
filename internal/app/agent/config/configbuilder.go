package config

import (
	"flag"
	"os"
)

type ConfigBuilder interface {
	SetAddress()
	SetPollInterval() error
	SetReportInterval() error
	GetConfig() Config
}

func GetConfigBuilder() ConfigBuilder {
	envAddress := os.Getenv("ADDRESS")
	envPoll := os.Getenv("POLL_INTERVAL")
	envReport := os.Getenv("REPORT_INTERVAL")
	if envAddress != "" && envPoll != "" && envReport != "" {
		return newEnvConfig()
	}
	cfg := newFlagsConfig()
	cfg.setFlagAddress()
	cfg.setFlagPoll()
	cfg.setFlagReport()
	flag.Parse()
	return cfg
}
