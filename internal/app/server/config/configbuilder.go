package config

import "os"

type ConfigBuilder interface {
	SetServerConfig()
	SetLoggerConfig()
	Build() Config
}

func GetConfigBuilder() ConfigBuilder {
	envEndpoint := os.Getenv("ADDRESS")
	envLogLevel := os.Getenv("LOG_LEVEL")

	if envEndpoint != "" && envLogLevel != "" {
		return newEnvConfig()

	}
	return newFlagsConfig()
}
