package config

import (
	"flag"
	"os"
)

type envConfig struct {
	sConfig     ServerConfig
	lConfig     LoggerConfig
	flagsValues flags
}

type flags struct {
	address  string
	logLevel string
}

func newEnvConfig() *envConfig {
	return &envConfig{}
}

func (ec *envConfig) setFlagsVariables() {
	flag.StringVar(&ec.flagsValues.address, "a", DefaultServerAddress, "address and port to run server")
	flag.StringVar(&ec.flagsValues.logLevel, "l", DefaultLogLevel, "log level")
	flag.Parse()
}

func (ec *envConfig) SetServerConfig() {

	envEndpoint := os.Getenv("ADDRESS")

	if envEndpoint != "" {
		ec.sConfig.Address = envEndpoint
		return
	}

	ec.sConfig.Address = ec.flagsValues.address
}

func (ec *envConfig) SetLoggerConfig() {

	envLogLevel := os.Getenv("LOG_LEVEL")

	if envLogLevel != "" {
		ec.lConfig.Level = envLogLevel
		return
	}

	ec.lConfig.Level = ec.flagsValues.logLevel
}

func (ec *envConfig) Build() Config {
	ec.setFlagsVariables()
	ec.SetServerConfig()
	ec.SetLoggerConfig()
	return Config{
		SConfig: ServerConfig{
			Address: ec.sConfig.Address,
		},
		LConfig: LoggerConfig{
			Level: ec.lConfig.Level,
		},
	}
}
