package config

import "time"

const (
	DefaultStoreInterval   = 300
	DefaultFileStoragePath = "storage.txt"
	DefaultRestore         = true
	MemStorageType         = "memory"
	FileStorageType        = "file"
	PostgresStorageType    = "postgres"
)

type DataBaseConfig struct {
	DSN             string
	FileStoragePath string
	DBType          string
	StoreInterval   time.Duration
	Restore         bool
}
