package config

import "flag"

type flagConfig struct {
	sConfig ServerConfig
	lConfig LoggerConfig
}

func newFlagsConfig() *flagConfig {
	return &flagConfig{}
}

func (fc *flagConfig) SetServerConfig() {
	flag.StringVar(&fc.sConfig.Address, "a", DefaultServerAddress, "address and port to run server")
}

func (fc *flagConfig) SetLoggerConfig() {
	flag.StringVar(&fc.lConfig.Level, "l", DefaultLogLevel, "log level")
}

func (fc *flagConfig) Build() Config {
	fc.SetServerConfig()
	fc.SetLoggerConfig()
	flag.Parse()
	return Config{
		SConfig: ServerConfig{
			Address: fc.sConfig.Address,
		},
		LConfig: LoggerConfig{
			Level: fc.lConfig.Level,
		},
	}
}
