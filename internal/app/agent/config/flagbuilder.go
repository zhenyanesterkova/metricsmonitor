package config

import (
	"errors"
	"flag"
	"strconv"
	"time"
)

func (fc *Config) readFlagAddress() {
	var temp string
	fc.addressFl = &temp
	flag.StringVar(fc.addressFl, "a", "localhost:8080", "address and port to run server")
}
func (fc *Config) readFlagPoll() {
	var temp int
	fc.pollIntervalFl = &temp
	flag.IntVar(fc.pollIntervalFl, "p", 2, "the frequency of polling metrics from the runtime package")
}
func (fc *Config) readFlagReport() {
	var temp int
	fc.reportIntervalFl = &temp
	flag.IntVar(fc.reportIntervalFl, "r", 10, "the frequency of sending metrics to the server")
}

func (fc *Config) SetFlagAddress() {
	if fc.addressFl != nil {
		fc.Address = *fc.addressFl
	}
}

func (fc *Config) SetFlagPollInterval() error {
	if fc.pollIntervalFl != nil {
		dur, err := time.ParseDuration(strconv.Itoa(*fc.pollIntervalFl) + "s")
		if err != nil {
			return errors.New("can not parse poll_interval as duration " + err.Error())
		}
		fc.PollInterval = dur
	}
	return nil
}

func (fc *Config) SetFlagReportInterval() error {
	if fc.reportIntervalFl != nil {
		dur, err := time.ParseDuration(strconv.Itoa(*fc.reportIntervalFl) + "s")
		if err != nil {
			return errors.New("can not parse report_interval as duration " + err.Error())
		}
		fc.ReportInterval = dur
	}
	return nil
}

func (fc *Config) BuildFlags() error {
	fc.readFlagAddress()
	fc.readFlagPoll()
	fc.readFlagReport()
	flag.Parse()

	fc.SetFlagAddress()

	err := fc.SetFlagPollInterval()
	if err != nil {
		return err
	}

	err = fc.SetFlagReportInterval()
	if err != nil {
		return err
	}

	return nil
}
