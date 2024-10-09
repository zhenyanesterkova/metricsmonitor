package config

import (
	"errors"
	"os"
	"time"
)

type envConfig struct {
	address        string
	pollInterval   time.Duration
	reportInterval time.Duration
}

func newEnvConfig() *envConfig {
	return &envConfig{}
}

func (ec *envConfig) SetAddress() {
	ec.address = os.Getenv("ADDRESS")
}

func (ec *envConfig) SetPollInterval() error {

	dur, err := time.ParseDuration(os.Getenv("POLL_INTERVAL") + "s")
	if err != nil {
		return errors.New("can not parse poll_interval as duration" + err.Error())
	}
	ec.pollInterval = dur

	return nil
}
func (ec *envConfig) SetReportInterval() error {

	dur, err := time.ParseDuration(os.Getenv("REPORT_INTERVAL") + "s")
	if err != nil {
		return errors.New("can not parse report_interval as duration" + err.Error())
	}
	ec.reportInterval = dur

	return nil
}

func (ec *envConfig) GetConfig() Config {
	return Config{
		Address:        ec.address,
		PollInterval:   ec.pollInterval,
		ReportInterval: ec.reportInterval,
	}
}
