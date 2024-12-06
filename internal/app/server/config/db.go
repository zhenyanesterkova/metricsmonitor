package config

import "time"

const (
	DefaultStoreInterval   = 300
	DefaultFileStoragePath = "storage.txt"
	DefaultRestore         = true
	DefaultMinRetryDelay   = time.Second
	DefaultMaxRetryDelay   = 5 * time.Second
	DefaultMaxRetryAttempt = 3
)

type DataBaseConfig struct {
	PostgresConfig    *PostgresConfig
	FileStorageConfig *FileStorageConfig
}

type PostgresConfig struct {
	DSN        string
	MinDelay   time.Duration
	MaxDelay   time.Duration
	MaxAttempt int
}

type FileStorageConfig struct {
	FileStoragePath string
	StoreInterval   time.Duration
	Restore         bool
}
