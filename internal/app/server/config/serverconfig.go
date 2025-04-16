package config

const (
	DefaultServerAddress = "localhost:8080"
	DefaultCryptoKeyPath = "../../../../build/crypto/private"
)

type ServerConfig struct {
	HashKey       *string
	Address       string
	CryptoKeyPath string
}
