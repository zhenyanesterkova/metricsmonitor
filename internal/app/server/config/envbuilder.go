package config

import "os"

type envConfig struct {
	sConfig ServerConfig
	lConfig LoggerConfig
}

func newEnvConfig() *envConfig {
	return &envConfig{}
}

func (ec *envConfig) SetServerConfig() {
	ec.sConfig.Address = os.Getenv("ADDRESS")
}
func (ec *envConfig) SetLoggerConfig() {
	ec.lConfig.Level = os.Getenv("LOG_LEVEL")
}

func (ec *envConfig) Build() Config {
	ec.SetServerConfig()
	ec.SetLoggerConfig()
	return Config{
		SConfig: ServerConfig{
			ec.sConfig.Address,
		},
		LConfig: LoggerConfig{
			ec.lConfig.Level,
		},
	}
}
