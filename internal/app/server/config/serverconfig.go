package config

const (
	DefaultServerAddress        = "localhost:8080"
	DefaultCryptoPrivateKeyPath = "./example-private.crt"
	DefaultCryptoPublicKeyPath  = "./example-public.crt"
)

type ServerConfig struct {
	HashKey              *string
	Address              string
	CryptoPrivateKeyPath string
	CryptoPublicKeyPath  string
}
