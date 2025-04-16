package config

const (
	DefaultServerAddress        = "localhost:8080"
	DefaultCryptoPrivateKeyPath = "../../build/crypto/private.crt"
	DefaultCryptoPublicKeyPath  = "../../build/crypto/public.crt"
)

type ServerConfig struct {
	HashKey              *string
	Address              string
	CryptoPrivateKeyPath string
	CryptoPublicKeyPath  string
}
