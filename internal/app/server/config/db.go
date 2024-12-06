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
	DSN string
}

type FileStorageConfig struct {
	FileStoragePath string
	StoreInterval   time.Duration
	Restore         bool
}
