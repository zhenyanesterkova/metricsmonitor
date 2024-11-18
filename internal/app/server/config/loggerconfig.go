package config

const (
	DefaultLogLevel = "info"
)

type LoggerConfig struct {
	levelParam *string
	Level      string
}
