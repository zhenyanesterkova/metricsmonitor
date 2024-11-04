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
	flag.StringVar(c.addressFl, "a", "localhost:8080", "address and port to run server")
}
func (c *Config) readFlagPoll() {
	var temp int
	c.pollIntervalFl = &temp
	flag.IntVar(c.pollIntervalFl, "p", 2, "the frequency of polling metrics from the runtime package")
}
func (c *Config) readFlagReport() {
	var temp int
	c.reportIntervalFl = &temp
	flag.IntVar(c.reportIntervalFl, "r", 10, "the frequency of sending metrics to the server")
}

func (c *Config) SetFlagAddress() {
	if c.addressFl != nil {
		c.Address = *c.addressFl
	}
}

func (c *Config) SetFlagPollInterval() error {
	if c.pollIntervalFl != nil {
		dur, err := time.ParseDuration(strconv.Itoa(*c.pollIntervalFl) + "s")
		if err != nil {
			return errors.New("can not parse poll_interval as duration " + err.Error())
		}
		c.PollInterval = dur
	}
	return nil
}

func (c *Config) SetFlagReportInterval() error {
	if c.reportIntervalFl != nil {
		dur, err := time.ParseDuration(strconv.Itoa(*c.reportIntervalFl) + "s")
		if err != nil {
			return errors.New("can not parse report_interval as duration " + err.Error())
		}
		c.ReportInterval = dur
	}
	return nil
}

func (c *Config) BuildFlags() error {
	c.readFlagAddress()
	c.readFlagPoll()
	c.readFlagReport()
	flag.Parse()

	c.SetFlagAddress()

	err := c.SetFlagPollInterval()
	if err != nil {
		return err
	}

	err = c.SetFlagReportInterval()
	if err != nil {
		return err
	}

	return nil
}
