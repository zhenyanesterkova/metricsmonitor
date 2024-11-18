package config

import "time"

const (
	DefaultStoreInterval   = 300
	DefaultFileStoragePath = "storage.txt"
	DefaultRestore         = true
)

type RestoreConfig struct {
	FileStoragePath string
	Restore         bool
	StoreInterval   time.Duration
}
