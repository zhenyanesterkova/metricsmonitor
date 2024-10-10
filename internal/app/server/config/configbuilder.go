package config

import "os"

type ConfigBuilder interface {
	SetServerConfig()
	GetConfig() Config
}

func GetConfigBuilder() ConfigBuilder {
	if envEndpoint := os.Getenv("ADDRESS"); envEndpoint != "" {
		return newEnvConfig()

	}
	return newFlagsConfig()
}