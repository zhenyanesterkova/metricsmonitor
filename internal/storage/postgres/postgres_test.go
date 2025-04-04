package postgres

import (
	"context"
	"errors"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
)

func TestPostgresStorage_Ping(t *testing.T) {
	pingErr := errors.New("error when ping DB")

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	pool.ExpectPing()

	psg := &PostgresStorage{
		log:  logger.LogrusLogger{},
		pool: pool,
	}

	t.Run("Success", func(t *testing.T) {
		err = psg.Ping()
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		pool.ExpectPing().WillReturnError(pingErr)
		err = psg.Ping()
		require.ErrorIs(t, err, pingErr)
	})
}

func TestPostgresStorage_UpdateMetric(t *testing.T) {
	counter := metric.New(metric.TypeCounter)
	counter.ID = "testCounter"
	*counter.Delta = 3

	gauge := metric.New(metric.TypeGauge)
	gauge.ID = "testGauge"
	*gauge.Value = 5.5

	tests := []struct {
		name    string
		psg     *PostgresStorage
		arg     metric.Metric
		wantErr bool
	}{
		{
			name:    "#Counter: succsess",
			arg:     counter,
			wantErr: false,
		},
		{
			name:    "#Counter: error",
			arg:     counter,
			wantErr: true,
		},
		{
			name:    "#Gauge: succsess",
			arg:     gauge,
			wantErr: false,
		},
		{
			name:    "#Gauge: error",
			arg:     gauge,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer pool.Close()

			var expQuery *pgxmock.ExpectedQuery
			switch tt.arg.MType {
			case metric.TypeCounter:
				expQuery = pool.ExpectQuery("INSERT INTO counters").
					WithArgs(tt.arg.ID, *tt.arg.Delta).
					WillReturnRows(
						pgxmock.NewRows([]string{"id", "cValue"}).
							AddRow(tt.arg.ID, *tt.arg.Delta),
					)
			case metric.TypeGauge:
				expQuery = pool.ExpectQuery("INSERT INTO gauges").
					WithArgs(tt.arg.ID, *tt.arg.Value).
					WillReturnRows(
						pgxmock.NewRows([]string{"id", "gValue"}).
							AddRow(tt.arg.ID, *tt.arg.Value),
					)
			}
			if tt.wantErr {
				expQuery.WillReturnError(errors.New("update failed"))
			}
			psg := PostgresStorage{
				log:  logger.LogrusLogger{},
				pool: pool,
			}
			got, err := psg.UpdateMetric(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostgresStorage.UpdateMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			assert.Equal(t, tt.arg, got)
		})
	}
}

func TestPostgresStorage_UpdateManyMetrics(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	psg := &PostgresStorage{
		log:  logger.NewLogrusLogger(),
		pool: pool,
	}

	counter := metric.New(metric.TypeCounter)
	counter.ID = "testCounter"
	*counter.Delta = 3

	gauge := metric.New(metric.TypeGauge)
	gauge.ID = "testGauge"
	*gauge.Value = 5.5

	t.Run("Success", func(t *testing.T) {
		arg := []metric.Metric{
			counter,
			gauge,
		}

		pool.ExpectBegin()
		pool.ExpectExec("INSERT INTO counters").
			WithArgs(arg[0].ID, *arg[0].Delta).
			WillReturnResult(pgxmock.NewResult("", 0))
		pool.ExpectExec("INSERT INTO gauges").
			WithArgs(arg[1].ID, *arg[1].Value).
			WillReturnResult(pgxmock.NewResult("", 0))
		pool.ExpectCommit()
		pool.ExpectRollback()

		psg.pool = pool

		err := psg.UpdateManyMetrics(context.TODO(), arg)
		if err != nil {
			t.Errorf("PostgresStorage.UpdateManyMetrics() error = %v", err)
		}
	})

	t.Run("Error: update counter", func(t *testing.T) {
		wantErr := errors.New("failed exec query update metric")
		arg := []metric.Metric{
			counter,
			gauge,
		}

		pool, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer pool.Close()

		pool.ExpectBegin()
		pool.ExpectExec("INSERT INTO counters").
			WithArgs(arg[0].ID, *arg[0].Delta).
			WillReturnError(wantErr)
		pool.ExpectRollback()

		psg.pool = pool

		err = psg.UpdateManyMetrics(context.TODO(), arg)
		if err == nil || !errors.Is(err, wantErr) {
			t.Errorf("PostgresStorage.UpdateManyMetrics(): error was expected")
		}
	})
	t.Run("Error: update gauge", func(t *testing.T) {
		wantErr := errors.New("failed exec query update metric")
		arg := []metric.Metric{
			counter,
			gauge,
		}

		pool, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer pool.Close()

		pool.ExpectBegin()
		pool.ExpectExec("INSERT INTO counters").
			WithArgs(arg[0].ID, *arg[0].Delta).
			WillReturnResult(pgxmock.NewResult("", 0))
		pool.ExpectExec("INSERT INTO gauges").
			WithArgs(arg[1].ID, *arg[1].Value).
			WillReturnError(wantErr)
		pool.ExpectRollback()

		psg.pool = pool

		err = psg.UpdateManyMetrics(context.TODO(), arg)
		if err == nil || !errors.Is(err, wantErr) {
			t.Errorf("PostgresStorage.UpdateManyMetrics(): error was expected")
		}
	})
	t.Run("Error: begin", func(t *testing.T) {
		wantErr := errors.New("failed start a transaction")
		arg := []metric.Metric{
			counter,
			gauge,
		}

		pool, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer pool.Close()

		pool.ExpectBegin().
			WillReturnError(wantErr)
		pool.ExpectRollback()

		psg.pool = pool

		err = psg.UpdateManyMetrics(context.TODO(), arg)
		if err == nil || !errors.Is(err, wantErr) {
			t.Errorf("PostgresStorage.UpdateManyMetrics(): error was expected")
		}
	})
	t.Run("Error: unknown metric type", func(t *testing.T) {
		unknown := metric.Metric{
			ID:    "unknown",
			MType: "unknown",
		}

		arg := []metric.Metric{
			counter,
			unknown,
		}

		pool.ExpectBegin()
		pool.ExpectExec("INSERT INTO counters").
			WithArgs(arg[0].ID, *arg[0].Delta).
			WillReturnResult(pgxmock.NewResult("", 0))
		pool.ExpectRollback()

		psg.pool = pool

		err := psg.UpdateManyMetrics(context.TODO(), arg)
		if err == nil || !errors.Is(err, ErrUnknownMetricType) {
			t.Errorf("PostgresStorage.UpdateManyMetrics() error = %v", err)
		}
	})
	t.Run("Error: tx.Commit", func(t *testing.T) {
		wantErr := errors.New("failed commits the transaction update metrics")
		arg := []metric.Metric{
			counter,
		}

		pool.ExpectBegin()
		pool.ExpectExec("INSERT INTO counters").
			WithArgs(arg[0].ID, *arg[0].Delta).
			WillReturnResult(pgxmock.NewResult("", 0))
		pool.ExpectCommit().
			WillReturnError(wantErr)

		psg.pool = pool

		err := psg.UpdateManyMetrics(context.TODO(), arg)
		if err == nil || !errors.Is(err, wantErr) {
			t.Errorf("PostgresStorage.UpdateManyMetrics() error = %v", err)
		}
	})
}

func TestPostgresStorage_GetAllMetrics(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	psg := &PostgresStorage{
		log:  logger.NewLogrusLogger(),
		pool: pool,
	}

	counter := metric.New(metric.TypeCounter)
	counter.ID = "testCounter"
	*counter.Delta = 3

	gauge := metric.New(metric.TypeGauge)
	gauge.ID = "testGauge"
	*gauge.Value = 5.5

	t.Run("Success", func(t *testing.T) {
		wantMetricList := [][2]string{
			{
				"testGauge", "5.5",
			},
			{
				"testCounter", "3",
			},
		}
		pool.ExpectQuery("SELECT id, g_value FROM gauges;").
			WillReturnRows(
				pgxmock.NewRows([]string{"id", "gVal"}).
					AddRow(gauge.ID, *gauge.Value),
			)
		pool.ExpectQuery("SELECT id, delta FROM counters;").
			WillReturnRows(
				pgxmock.NewRows([]string{"id", "cVal"}).
					AddRow(counter.ID, *counter.Delta),
			)

		psg.pool = pool

		metricList, err := psg.GetAllMetrics()
		require.Equal(t, wantMetricList, metricList)
		if err != nil {
			t.Errorf("PostgresStorage.UpdateManyMetrics() error = %v", err)
		}
	})
	t.Run("Error: get gauges", func(t *testing.T) {
		wantErr := errors.New("failed to select all metrics from gauges table")

		pool, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer pool.Close()
		pool.ExpectQuery("SELECT id, g_value FROM gauges;").
			WillReturnError(wantErr)

		psg.pool = pool

		_, err = psg.GetAllMetrics()
		if err == nil || !errors.Is(err, wantErr) {
			t.Errorf("PostgresStorage.UpdateManyMetrics() error = %v", err)
		}
	})
	t.Run("Error: get counters", func(t *testing.T) {
		wantErr := errors.New("failed to select all metrics from counters table")

		pool, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer pool.Close()
		pool.ExpectQuery("SELECT id, g_value FROM gauges;").
			WillReturnRows(
				pgxmock.NewRows([]string{"id", "gVal"}).
					AddRow(gauge.ID, *gauge.Value),
			)
		pool.ExpectQuery("SELECT id, delta FROM counters;").
			WillReturnError(wantErr)

		psg.pool = pool

		_, err = psg.GetAllMetrics()
		if err == nil || !errors.Is(err, wantErr) {
			t.Errorf("PostgresStorage.UpdateManyMetrics() error = %v", err)
		}
	})
	t.Run("Error: scan gauges", func(t *testing.T) {
		pool, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer pool.Close()
		pool.ExpectQuery("SELECT id, g_value FROM gauges;").
			WillReturnRows(
				pgxmock.NewRows([]string{"id"}).
					AddRow(gauge.ID),
			)

		psg.pool = pool

		_, err = psg.GetAllMetrics()
		if err == nil || !errors.Is(err, ErrScanGauges) {
			t.Errorf("PostgresStorage.UpdateManyMetrics() error = %v", err)
		}
	})
	t.Run("Error: scan counters", func(t *testing.T) {
		pool, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer pool.Close()
		pool.ExpectQuery("SELECT id, g_value FROM gauges;").
			WillReturnRows(
				pgxmock.NewRows([]string{"id", "gVal"}).
					AddRow(gauge.ID, *gauge.Value),
			)
		pool.ExpectQuery("SELECT id, delta FROM counters;").
			WillReturnRows(
				pgxmock.NewRows([]string{"id"}).
					AddRow(counter.ID),
			)

		psg.pool = pool

		_, err = psg.GetAllMetrics()
		if err == nil || !errors.Is(err, ErrScanCounters) {
			t.Errorf("PostgresStorage.UpdateManyMetrics() error = %v", err)
		}
	})
}

func TestPostgresStorage_GetMetricValue(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	psg := &PostgresStorage{
		log:  logger.NewLogrusLogger(),
		pool: pool,
	}

	counter := metric.New(metric.TypeCounter)
	counter.ID = "testCounter"
	*counter.Delta = 3

	gauge := metric.New(metric.TypeGauge)
	gauge.ID = "testGauge"
	*gauge.Value = 5.5

	unknownGauge := metric.New(metric.TypeGauge)
	unknownGauge.ID = "unknownTestGauge"
	*unknownGauge.Value = 5.5

	unknownCounter := metric.New(metric.TypeCounter)
	unknownCounter.ID = "unknownTestCounter"
	*unknownCounter.Delta = 3

	t.Run("Success: get counter", func(t *testing.T) {
		pool.ExpectQuery("SELECT id, delta FROM counters").
			WithArgs(counter.ID).
			WillReturnRows(
				pgxmock.NewRows([]string{"id", "cValue"}).
					AddRow(counter.ID, *counter.Delta),
			)

		psg.pool = pool

		got, err := psg.GetMetricValue(counter.ID, counter.MType)
		if err != nil {
			t.Errorf("PostgresStorage.GetMetricValue() error = %v", err)
			return
		}
		require.Equal(t, counter, got)
	})
	t.Run("Success: get gauge", func(t *testing.T) {
		pool.ExpectQuery("SELECT id, g_value FROM gauges").
			WithArgs(gauge.ID).
			WillReturnRows(
				pgxmock.NewRows([]string{"id", "gValue"}).
					AddRow(gauge.ID, *gauge.Value),
			)

		psg.pool = pool

		got, err := psg.GetMetricValue(gauge.ID, gauge.MType)
		if err != nil {
			t.Errorf("PostgresStorage.GetMetricValue() error = %v", err)
			return
		}
		require.Equal(t, gauge, got)
	})
	t.Run("Error: get unknown counter", func(t *testing.T) {
		pool.ExpectQuery("SELECT id, delta FROM counters").
			WithArgs(unknownCounter.ID).
			WillReturnRows(pgxmock.NewRows([]string{}))

		psg.pool = pool

		_, err := psg.GetMetricValue(unknownCounter.ID, unknownCounter.MType)
		if err == nil || !errors.Is(err, metric.ErrUnknownMetric) {
			t.Errorf("PostgresStorage.GetMetricValue() error = %v", err)
			return
		}
	})
	t.Run("Error: when scan counter", func(t *testing.T) {
		pool.ExpectQuery("SELECT id, delta FROM counters").
			WithArgs(counter.ID).
			WillReturnRows(
				pgxmock.NewRows([]string{"id"}).
					AddRow(counter.ID),
			)

		psg.pool = pool

		_, err := psg.GetMetricValue(counter.ID, counter.MType)
		psg.log.LogrusLog.Infof("Err: %v", err)
		if err == nil || !errors.Is(err, ErrScanCounter) {
			t.Errorf("PostgresStorage.GetMetricValue() error = %v", err)
			return
		}
	})
	t.Run("Error: get unknown gauge", func(t *testing.T) {
		pool.ExpectQuery("SELECT id, g_value FROM gauges").
			WithArgs(unknownGauge.ID).
			WillReturnRows(pgxmock.NewRows([]string{}))

		psg.pool = pool

		_, err := psg.GetMetricValue(unknownGauge.ID, unknownGauge.MType)
		if err == nil || !errors.Is(err, metric.ErrUnknownMetric) {
			t.Errorf("PostgresStorage.GetMetricValue() error = %v", err)
			return
		}
	})
	t.Run("Error: when scan gauge", func(t *testing.T) {
		pool.ExpectQuery("SELECT id, g_value FROM gauges").
			WithArgs(gauge.ID).
			WillReturnRows(
				pgxmock.NewRows([]string{"id"}).
					AddRow(gauge.ID),
			)

		psg.pool = pool

		_, err := psg.GetMetricValue(gauge.ID, gauge.MType)
		psg.log.LogrusLog.Infof("Err: %v", err)
		if err == nil || !errors.Is(err, ErrScanGauge) {
			t.Errorf("PostgresStorage.GetMetricValue() error = %v", err)
			return
		}
	})
}

func TestPostgresStorage_Close(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	psg := &PostgresStorage{
		log:  logger.NewLogrusLogger(),
		pool: pool,
	}

	t.Run("Success", func(t *testing.T) {
		if err := psg.Close(); err != nil {
			t.Errorf("PostgresStorage.Close() error = %v", err)
		}
	})
}
