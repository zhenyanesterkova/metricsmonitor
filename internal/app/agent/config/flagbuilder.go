package config

import (
	"flag"
)

func (c *Config) readFlagAddress() {
	flag.StringVar(&c.Address, "a", c.Address, "address and port to run server")
}
func (c *Config) readFlagPoll() {
	flag.DurationVar(&c.PollInterval, "p", c.PollInterval, "the frequency of polling metrics from the runtime package")
}
func (c *Config) readFlagReport() {
	flag.DurationVar(&c.ReportInterval, "r", c.ReportInterval, "the frequency of sending metrics to the server")
}

func (c *Config) buildFlags() {
	c.readFlagAddress()
	c.readFlagPoll()
	c.readFlagReport()
	flag.Parse()
}
