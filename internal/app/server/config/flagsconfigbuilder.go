package config

import (
	"flag"
)

func (c *Config) setFlagsVariables() {
	flag.StringVar(&c.SConfig.Address, "a", c.SConfig.Address, "address and port to run server")
	flag.StringVar(&c.LConfig.Level, "l", c.LConfig.Level, "log level")
	flag.DurationVar(&c.RConfig.StoreInterval, "i", c.RConfig.StoreInterval, "store interval")
	flag.StringVar(&c.RConfig.FileStoragePath, "f", c.RConfig.FileStoragePath, "file storage path")
	flag.BoolVar(&c.RConfig.Restore, "r", c.RConfig.Restore, "need restore")
	flag.Parse()
}

func (c *Config) flagBuild() {
	c.setFlagsVariables()
}
