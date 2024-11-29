package config

import (
	"time"
)

const (
	DefaultStoreInterval   = 300
	DefaultFileStoragePath = "storage.txt"
	DefaultRestore         = true
)

type RestoreConfig struct {
	FileStoragePath string
	StoreInterval   time.Duration
	Restore         bool
}
