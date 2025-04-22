package config

import "time"

const (
	DefaultStoreInterval   = 300
	DefaultFileStoragePath = "storage.txt"
	DefaultRestore         = true
)

type DataBaseConfig struct {
	PostgresConfig    *PostgresConfig    `json:"postgres_storage"`
	FileStorageConfig *FileStorageConfig `json:"file_storage"`
}

type PostgresConfig struct {
	DSN string `json:"database_dsn"`
}

type FileStorageConfig struct {
	FileStoragePath string        `json:"store_file"`
	StoreInterval   time.Duration `json:"store_interval"`
	Restore         bool          `json:"restore"`
}
