package config

const (
	DefaultServerAddress        = "localhost:8080"
	DefaultCryptoPrivateKeyPath = "../../build/private.crt"
	DefaultCryptoPublicKeyPath  = "../../build/public.crt"
)

type ServerConfig struct {
	HashKey              *string
	Address              string
	CryptoPrivateKeyPath string
	CryptoPublicKeyPath  string
}
