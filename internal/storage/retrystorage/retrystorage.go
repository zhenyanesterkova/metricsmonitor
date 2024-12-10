package retrystorage

import (
	"context"
	"fmt"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/backoff"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage"
)

type RetryStorage struct {
	storage    storage.Store
	backoff    *backoff.Backoff
	logger     logger.LogrusLogger
	checkRetry func(error) bool
}

func New(
	cfg config.DataBaseConfig,
	loggerInst logger.LogrusLogger,
	bf *backoff.Backoff,
	checkRetryFunc func(error) bool,
) (
	*RetryStorage,
	error,
) {
	retryStore := &RetryStorage{
		checkRetry: checkRetryFunc,
		backoff:    bf,
		logger:     loggerInst,
	}

	store, err := storage.NewStore(cfg, loggerInst)
	if err != nil {
		if retryStore.checkRetry(err) {
			err = retryStore.retry(func() error {
				store, err = storage.NewStore(cfg, loggerInst)
				if err != nil {
					return fmt.Errorf("failed retry create storage: %w", err)
				}
				return nil
			})
		}
		loggerInst.LogrusLog.Errorf("can not create storage: %v", err)
		return retryStore, fmt.Errorf("can not create storage: %w", err)
	}

	retryStore.storage = store
	return retryStore, nil
}

func (rs *RetryStorage) UpdateMetric(m metric.Metric) (metric.Metric, error) {
	resMetric, err := rs.storage.UpdateMetric(m)
	if rs.checkRetry(err) {
		err = rs.retry(func() error {
			resMetric, err = rs.storage.UpdateMetric(m)
			if err != nil {
				return fmt.Errorf("failed retry update metric: %w", err)
			}
			return nil
		})
	}
	if err != nil {
		return resMetric, fmt.Errorf("failed update: %w", err)
	}
	return resMetric, nil
}

func (rs *RetryStorage) GetAllMetrics() ([][2]string, error) {
	resMetricList, err := rs.storage.GetAllMetrics()
	if rs.checkRetry(err) {
		err = rs.retry(func() error {
			resMetricList, err = rs.storage.GetAllMetrics()
			if err != nil {
				return fmt.Errorf("failed retry get metrics list: %w", err)
			}
			return nil
		})
	}
	if err != nil {
		return resMetricList, fmt.Errorf("failed get metrics: %w", err)
	}
	return resMetricList, nil
}

func (rs *RetryStorage) GetMetricValue(name, typeMetric string) (metric.Metric, error) {
	resMetric, err := rs.storage.GetMetricValue(name, typeMetric)
	if rs.checkRetry(err) {
		err = rs.retry(func() error {
			resMetric, err = rs.storage.GetMetricValue(name, typeMetric)
			if err != nil {
				return fmt.Errorf("failed retry get metric: %w", err)
			}
			return nil
		})
	}
	if err != nil {
		return resMetric, fmt.Errorf("failed get metric: %w", err)
	}
	return resMetric, nil
}

func (rs *RetryStorage) UpdateManyMetrics(ctx context.Context, mList []metric.Metric) error {
	err := rs.storage.UpdateManyMetrics(ctx, mList)
	if rs.checkRetry(err) {
		err = rs.retry(func() error {
			err = rs.storage.UpdateManyMetrics(ctx, mList)
			if err != nil {
				return fmt.Errorf("failed retry update metrics: %w", err)
			}
			return nil
		})
	}
	if err != nil {
		return fmt.Errorf("failed update metrics: %w", err)
	}
	return nil
}

func (rs *RetryStorage) Ping() error {
	err := rs.storage.Ping()
	if rs.checkRetry(err) {
		err = rs.retry(func() error {
			err = rs.storage.Ping()
			if err != nil {
				return fmt.Errorf("failed retry ping: %w", err)
			}
			return nil
		})
	}
	if err != nil {
		return fmt.Errorf("failed ping: %w", err)
	}
	return nil
}

func (rs *RetryStorage) Close() error {
	if err := rs.storage.Close(); err != nil {
		return fmt.Errorf("failed close DB: %w", err)
	}
	return nil
}

func (rs *RetryStorage) retry(work func() error) error {
	log := rs.logger.LogrusLog
	defer rs.backoff.Reset()
	for {
		log.Debug("attempt to repeat ...")
		err := work()

		if err == nil {
			return nil
		}

		if rs.checkRetry(err) {
			var delay time.Duration
			if delay = rs.backoff.Next(); delay == backoff.Stop {
				return err
			}
			time.Sleep(delay)
		}
	}
}
