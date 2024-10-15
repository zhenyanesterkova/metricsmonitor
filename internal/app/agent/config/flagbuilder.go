package config

import (
	"errors"
	"flag"
	"strconv"
	"time"
)

type flagConfig struct {
	address          string
	pollInterval     time.Duration
	reportInterval   time.Duration
	parseFlagsStruct *flags
}

type flags struct {
	addressFl string
	pollFl    int
	reportFl  int
}

func newFlagsConfig() *flagConfig {
	return &flagConfig{
		parseFlagsStruct: &flags{},
	}
}

func (fc *flagConfig) setFlagAddress() {
	flag.StringVar(&fc.parseFlagsStruct.addressFl, "a", "localhost:8080", "address and port to run server")
}
func (fc *flagConfig) setFlagPoll() {
	flag.IntVar(&fc.parseFlagsStruct.pollFl, "p", 2, "the frequency of polling metrics from the runtime package")
}
func (fc *flagConfig) setFlagReport() {
	flag.IntVar(&fc.parseFlagsStruct.reportFl, "r", 10, "the frequency of sending metrics to the server")
}

func (fc *flagConfig) SetAddress() {
	fc.address = fc.parseFlagsStruct.addressFl
}
func (fc *flagConfig) SetPollInterval() error {

	dur, err := time.ParseDuration(strconv.Itoa(fc.parseFlagsStruct.pollFl) + "s")
	if err != nil {
		return errors.New("can not parse poll_interval as duration " + err.Error())
	}
	fc.pollInterval = dur

	return nil
}
func (fc *flagConfig) SetReportInterval() error {

	dur, err := time.ParseDuration(strconv.Itoa(fc.parseFlagsStruct.reportFl) + "s")
	if err != nil {
		return errors.New("can not parse report_interval as duration " + err.Error())
	}
	fc.reportInterval = dur

	return nil
}

func (fc *flagConfig) Build() (Config, error) {
	fc.SetAddress()

	err := fc.SetPollInterval()
	if err != nil {
		return Config{}, err
	}

	err = fc.SetReportInterval()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Address:        fc.address,
		PollInterval:   fc.pollInterval,
		ReportInterval: fc.reportInterval,
	}, nil
}
