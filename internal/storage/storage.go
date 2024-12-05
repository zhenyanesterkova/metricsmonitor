package storage

import (
	"context"
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
	UpdateManyMetrics(ctx context.Context, mList []metric.Metric) error
}

func NewStore(conf config.DataBaseConfig, log logger.LogrusLogger) (Store, error) {
	if conf.PostgresConfig != nil {
		store, err := postgres.New(conf.PostgresConfig.DSN, log)
		if err != nil {
			return nil, fmt.Errorf("failed create postgres storage: %w", err)
		}
		return store, nil
	}

	if conf.FileStorageConfig != nil {
		store, err := filestorage.New(*conf.FileStorageConfig, log)
		if err != nil {
			return nil, fmt.Errorf("failed create file storage: %w", err)
		}
		return store, nil
	}

	return memstorage.New(), nil
}
