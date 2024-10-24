package config

type ConfigBuilder interface {
	SetServerConfig()
	SetLoggerConfig()
	Build() Config
}

func GetConfigBuilder() ConfigBuilder {
	return newEnvConfig()
}
