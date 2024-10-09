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

func (d *ConfigDirector) BuildConfig() (Config, error) {
	d.Builder.SetAddress()

	err := d.Builder.SetPollInterval()
	if err != nil {
		return d.Builder.GetConfig(), err
	}

	d.Builder.SetReportInterval()
	if err != nil {
		return d.Builder.GetConfig(), err
	}
	return d.Builder.GetConfig(), nil
}
