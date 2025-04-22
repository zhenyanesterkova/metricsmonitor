package config

const (
	DefaultServerAddress        = "localhost:8080"
	DefaultCryptoPrivateKeyPath = "example-private.crt"
	DefaultCryptoPublicKeyPath  = "example-public.crt"
	DefualtNeedGenKeys          = false
	DefaultConfigsFileName      = "config.json"
)

type ServerConfig struct {
	HashKey              *string `json:"hashkey"`
	Address              string  `json:"address"`
	CryptoPrivateKeyPath string  `json:"crypto_key"`
	CryptoPublicKeyPath  string  `json:"crypto_pub_key"`
	ConfigsFileName      string  `json:"config_file_name"`
}
