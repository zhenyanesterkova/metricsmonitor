package config

const (
	DefaultServerAddress = "localhost:8080"
)

type ServerConfig struct {
	addressParam *string
	Address      string
}
