package config

import (
	"errors"
	"flag"
	"strconv"
	"time"
)

func (c *Config) readFlagAddress() {
	var temp string
	c.addressFl = &temp
	flag.StringVar(c.addressFl, "a", defaultAddress, "address and port to run server")
}
func (c *Config) readFlagPoll() {
	var temp int
	c.pollIntervalFl = &temp
	flag.IntVar(c.pollIntervalFl, "p", defaultPollInt, "the frequency of polling metrics from the runtime package")
}
func (c *Config) readFlagReport() {
	var temp int
	c.reportIntervalFl = &temp
	flag.IntVar(c.reportIntervalFl, "r", defaultReportInt, "the frequency of sending metrics to the server")
}

func (c *Config) setFlagAddress() {
	if c.addressFl != nil {
		c.Address = *c.addressFl
	}
}

func (c *Config) setFlagPollInterval() error {
	if c.pollIntervalFl != nil {
		dur, err := time.ParseDuration(strconv.Itoa(*c.pollIntervalFl) + "s")
		if err != nil {
			return errors.New("can not parse poll_interval as duration " + err.Error())
		}
		c.PollInterval = dur
	}
	return nil
}

func (c *Config) setFlagReportInterval() error {
	if c.reportIntervalFl != nil {
		dur, err := time.ParseDuration(strconv.Itoa(*c.reportIntervalFl) + "s")
		if err != nil {
			return errors.New("can not parse report_interval as duration " + err.Error())
		}
		c.ReportInterval = dur
	}
	return nil
}

func (c *Config) buildFlags() error {
	c.readFlagAddress()
	c.readFlagPoll()
	c.readFlagReport()
	flag.Parse()

	c.setFlagAddress()

	err := c.setFlagPollInterval()
	if err != nil {
		return err
	}

	err = c.setFlagReportInterval()
	if err != nil {
		return err
	}

	return nil
}
