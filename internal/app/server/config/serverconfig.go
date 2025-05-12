package config

import "net"

const (
	DefaultServerAddress        = "localhost:8080"
	DefaultCryptoPrivateKeyPath = "example-private.crt"
	DefaultCryptoPublicKeyPath  = "example-public.crt"
	DefualtNeedGenKeys          = false
	DefaultConfigsFileName      = "server_config.json"
)

type ServerConfig struct {
	HashKey              *string    `json:"hashkey"`
	TIpNet               *net.IPNet `json:"-"`
	Address              string     `json:"address"`
	CryptoPrivateKeyPath string     `json:"crypto_key"`
	CryptoPublicKeyPath  string     `json:"crypto_pub_key"`
	ConfigsFileName      string     `json:"config_file_name"`
	StringCIDR           string     `json:"trusted_subnet"`
}
