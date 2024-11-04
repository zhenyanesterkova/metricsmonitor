package config

import (
	"errors"
	"os"
	"time"
)

func (ec *Config) ReadEnv() {
	addr := os.Getenv("ADDRESS")
	pollInt := os.Getenv("POLL_INTERVAL")
	reportInt := os.Getenv("REPORT_INTERVAL")
	if addr != "" {
		ec.addressEnv = &addr
	}
	if pollInt != "" {
		ec.pollIntervalEnv = &pollInt
	}
	if reportInt != "" {
		ec.reportIntervalEnv = &reportInt
	}
}

func (ec *Config) SetEnvAddress() {
	if ec.addressEnv != nil {
		ec.Address = *ec.addressEnv
	}
}

func (ec *Config) SetEnvPollInterval() error {
	if ec.pollIntervalEnv != nil {
		dur, err := time.ParseDuration(*ec.pollIntervalEnv + "s")
		if err != nil {
			return errors.New("can not parse poll_interval as duration" + err.Error())
		}
		ec.PollInterval = dur
	}
	return nil
}
func (ec *Config) SetEnvReportInterval() error {
	if ec.reportIntervalEnv != nil {
		dur, err := time.ParseDuration(*ec.reportIntervalEnv + "s")
		if err != nil {
			return errors.New("can not parse report_interval as duration" + err.Error())
		}
		ec.ReportInterval = dur
	}
	return nil
}

func (ec *Config) BuildEnv() error {
	ec.ReadEnv()

	ec.SetEnvAddress()

	err := ec.SetEnvPollInterval()
	if err != nil {
		return err
	}

	err = ec.SetEnvReportInterval()
	if err != nil {
		return err
	}

	return nil
}
