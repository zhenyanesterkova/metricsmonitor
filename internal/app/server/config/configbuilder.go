package config

type ConfigBuilder interface {
	SetServerConfig()
	SetLoggerConfig()
	SetRestoreConfig() error
	Build() (Config, error)
}

func GetConfigBuilder() ConfigBuilder {
	return newEnvConfig()
}
