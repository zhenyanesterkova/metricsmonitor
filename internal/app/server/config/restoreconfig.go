package config

import "time"

const (
	DefaultStoreInterval   = 300
	DefaultFileStoragePath = "/storage/store.txt"
	DefaultRestore         = true
)

type RestoreConfig struct {
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
}
