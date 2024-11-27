package storage

import (
	"fmt"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/filestorage"
)

type Store interface {
	UpdateMetric(metric.Metric) (metric.Metric, error)
	GetAllMetrics() ([][2]string, error)
	GetMetricValue(name, typeMetric string) (metric.Metric, error)
	Close() error
}

func NewStore(conf config.RestoreConfig, log logger.LogrusLogger) (Store, error) {
	store, err := filestorage.New(conf, log)
	if err != nil {
		return nil, fmt.Errorf("create storage error: %w", err)
	}
	return store, nil
}
