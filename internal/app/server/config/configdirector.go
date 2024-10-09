package config

type ConfigDirector struct {
	Builder ConfigBuilder
}

func NewConfigDirector(b ConfigBuilder) *ConfigDirector {
	return &ConfigDirector{
		Builder: b,
	}
}

func (d *ConfigDirector) SetConfigBuilder(b ConfigBuilder) {
	d.Builder = b
}

func (d *ConfigDirector) BuildConfig() Config {
	d.Builder.SetServerConfig()
	return d.Builder.GetConfig()
}
