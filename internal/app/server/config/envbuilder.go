package config

import "os"

type envConfig struct {
	sConfig ServerConfig
}

func newEnvConfig() *envConfig {
	return &envConfig{}
}

func (ec *envConfig) SetServerConfig() {
	ec.sConfig.Address = os.Getenv("ADDRESS")
}

func (ec *envConfig) GetConfig() Config {
	return Config{
		SConfig: ServerConfig{
			ec.sConfig.Address,
		},
	}
}
