package config

import "flag"

type flagConfig struct {
	sConfig ServerConfig
}

func newFlagsConfig() *flagConfig {
	return &flagConfig{}
}

func (fc *flagConfig) SetServerConfig() {
	flag.StringVar(&fc.sConfig.Address, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
}

func (fc *flagConfig) GetConfig() Config {
	return Config{
		SConfig: ServerConfig{
			Address: fc.sConfig.Address,
		},
	}
}
