package config

import "time"

const (
	DefaultStoreInterval   = 300
	DefaultFileStoragePath = "storage.txt"
	DefaultRestore         = true
)

type DataBaseConfig struct {
	PostgresConfig    *PostgresConfig
	FileStorageConfig *FileStorageConfig
}

type PostgresConfig struct {
	DSN string `json:"database_dsn"`
}

type FileStorageConfig struct {
	FileStoragePath string        `json:"store_file"`
	StoreInterval   time.Duration `json:"store_interval"`
	Restore         bool          `json:"restore"`
}
