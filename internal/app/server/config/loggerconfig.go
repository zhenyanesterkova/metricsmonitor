package config

const (
	DefaultLogLevel = "info"
)

type LoggerConfig struct {
	Level string `json:"log_level"`
}
