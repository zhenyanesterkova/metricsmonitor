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
	filePathParam        *string
	restoreParam         *bool
	storeIntervalParamFl *int
	FileStoragePath      string
	StoreInterval        time.Duration
	Restore              bool
}
