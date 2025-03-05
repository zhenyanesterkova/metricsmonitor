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
		err = psg.pool.Ping(context.TODO())
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		pool.ExpectPing().WillReturnError(pingErr)
		err = psg.pool.Ping(context.TODO())
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
