package storage

import (
	"errors"
	"fmt"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/filestorage"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memstorage"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/postgres"
)

type Store interface {
	UpdateMetric(metric.Metric) (metric.Metric, error)
	GetAllMetrics() ([][2]string, error)
	GetMetricValue(name, typeMetric string) (metric.Metric, error)
	Close() error
	Ping() (bool, error)
}

func NewStore(conf config.DataBaseConfig, log logger.LogrusLogger) (Store, error) {
	switch conf.DBType {
	case config.PostgresStorageType:
		store, err := postgres.New(conf.DSN)
		if err != nil {
			return nil, fmt.Errorf("failed create postgres storage: %w", err)
		}
		return store, nil
	case config.MemStorageType:
		return memstorage.New(), nil
	case config.FileStorageType:
		store, err := filestorage.New(conf, log)
		if err != nil {
			return nil, fmt.Errorf("failed create file storage: %w", err)
		}
		return store, nil
	default:
		return nil, errors.New("failed create storage: unknown storage type")
	}
}
