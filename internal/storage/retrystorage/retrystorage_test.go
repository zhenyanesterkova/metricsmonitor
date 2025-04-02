package retrystorage

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/backoff"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
)

type mockStore struct {
	count   int
	success bool
}

func (m *mockStore) UpdateMetric(metric.Metric) (metric.Metric, error) {
	pgErr := &pgconn.PgError{}
	pgErr.Code = pgerrcode.ConnectionException
	if !m.success || m.count == 0 {
		m.count++
		return metric.Metric{}, pgErr
	}

	m.count++
	return metric.New("counter"), nil
}

func (m *mockStore) GetAllMetrics() ([][2]string, error) {
	pgErr := &pgconn.PgError{}
	pgErr.Code = pgerrcode.ConnectionException
	if !m.success || m.count == 0 {
		m.count++
		return [][2]string{}, pgErr
	}

	m.count++
	return [][2]string{}, nil
}

func (m *mockStore) GetMetricValue(name, typeMetric string) (metric.Metric, error) {
	pgErr := &pgconn.PgError{}
	pgErr.Code = pgerrcode.ConnectionException
	if !m.success || m.count == 0 {
		m.count++
		return metric.Metric{}, pgErr
	}

	m.count++
	return metric.Metric{}, nil
}

func (m *mockStore) Close() error {
	m.count++
	return errors.New("err close")
}

func (m *mockStore) Ping() error {
	pgErr := &pgconn.PgError{}
	pgErr.Code = pgerrcode.ConnectionException
	if !m.success || m.count == 0 {
		m.count++
		return pgErr
	}

	m.count++
	return nil
}

func (m *mockStore) UpdateManyMetrics(ctx context.Context, mList []metric.Metric) error {
	pgErr := &pgconn.PgError{}
	pgErr.Code = pgerrcode.ConnectionException
	if !m.success || m.count == 0 {
		m.count++
		return pgErr
	}

	m.count++
	return nil
}

func TestRetryStorage(t *testing.T) {
	mock := &mockStore{
		count:   0,
		success: true,
	}

	log := logger.NewLogrusLogger()

	backoffInst := backoff.New(
		5,
		10,
		config.DefaultMaxRetryAttempt,
	)

	checkRetryFunc := func(err error) bool {
		var pgErr *pgconn.PgError
		var pgErrConn *pgconn.ConnectError
		res := false
		if errors.As(err, &pgErr) {
			res = pgerrcode.IsConnectionException(pgErr.Code)
		} else if errors.As(err, &pgErrConn) {
			res = true
		}
		return res
	}

	retryStore := RetryStorage{
		storage:    mock,
		logger:     log,
		backoff:    backoffInst,
		checkRetry: checkRetryFunc,
	}

	t.Run("success second try ping", func(t *testing.T) {
		err := retryStore.Ping()
		require.NoError(t, err)
	})
	t.Run("failed retry ping", func(t *testing.T) {
		mock.success = false
		retryStore.storage = mock

		err := retryStore.Ping()
		require.Error(t, err)
	})

	t.Run("success second try UpdateMetric", func(t *testing.T) {
		mock.count = 0
		mock.success = true
		retryStore.storage = mock
		metrica := metric.New("counter")
		_, err := retryStore.UpdateMetric(metrica)
		require.NoError(t, err)
	})
	t.Run("failed retry UpdateMetric", func(t *testing.T) {
		mock.success = false
		retryStore.storage = mock
		metrica := metric.New("counter")
		_, err := retryStore.UpdateMetric(metrica)
		require.Error(t, err)
	})

	t.Run("success second try GetAllMetrics", func(t *testing.T) {
		mock.count = 0
		mock.success = true
		retryStore.storage = mock
		_, err := retryStore.GetAllMetrics()
		require.NoError(t, err)
	})
	t.Run("failed retry GetAllMetrics", func(t *testing.T) {
		mock.success = false
		retryStore.storage = mock
		_, err := retryStore.GetAllMetrics()
		require.Error(t, err)
	})

	t.Run("success second try GetMetricValue", func(t *testing.T) {
		mock.count = 0
		mock.success = true
		retryStore.storage = mock
		_, err := retryStore.GetMetricValue("", "")
		require.NoError(t, err)
	})
	t.Run("failed retry GetMetricValue", func(t *testing.T) {
		mock.success = false
		retryStore.storage = mock
		_, err := retryStore.GetMetricValue("", "")
		require.Error(t, err)
	})

	t.Run("success second try UpdateManyMetrics", func(t *testing.T) {
		mock.count = 0
		mock.success = true
		retryStore.storage = mock
		err := retryStore.UpdateManyMetrics(context.TODO(), []metric.Metric{})
		require.NoError(t, err)
	})
	t.Run("failed retry UpdateManyMetrics", func(t *testing.T) {
		mock.success = false
		retryStore.storage = mock
		err := retryStore.UpdateManyMetrics(context.TODO(), []metric.Metric{})
		require.Error(t, err)
	})

	t.Run("failed retry Close", func(t *testing.T) {
		err := retryStore.Close()
		require.Error(t, err)
	})
}
