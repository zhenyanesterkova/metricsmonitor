package config

const (
	DefaultServerAddress        = "localhost:8080"
	DefaultCryptoPrivateKeyPath = "../../build/example-private.crt"
	DefaultCryptoPublicKeyPath  = "../../build/example-public.crt"
)

type ServerConfig struct {
	HashKey              *string
	Address              string
	CryptoPrivateKeyPath string
	CryptoPublicKeyPath  string
}
