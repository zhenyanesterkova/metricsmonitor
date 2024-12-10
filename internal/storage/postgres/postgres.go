package postgres

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/storagerror"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
	log  logger.LogrusLogger
}

func New(dsn string, lg logger.LogrusLogger) (*PostgresStorage, error) {
	if err := runMigrations(dsn); err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create a connection pool: %w", err)
	}
	return &PostgresStorage{
		pool: pool,
		log:  lg,
	}, nil
}

//go:embed migrations/*.sql
var migrationsDir embed.FS

func runMigrations(dsn string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		if checkRetry(err) {
			err = storagerror.NewRetriableError(err)
		}
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			if checkRetry(err) {
				err = storagerror.NewRetriableError(err)
			}
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
	}
	return nil
}

func (psg *PostgresStorage) Ping() error {
	if err := psg.pool.Ping(context.TODO()); err != nil {
		if checkRetry(err) {
			err = storagerror.NewRetriableError(err)
		}
		return fmt.Errorf("failed to ping the DB: %w", err)
	}

	return nil
}

func (psg *PostgresStorage) UpdateMetric(m metric.Metric) (metric.Metric, error) {
	var updating metric.Metric
	var row pgx.Row
	var id string
	if m.MType == metric.TypeGauge {
		var gValue float64
		row = psg.pool.QueryRow(
			context.TODO(),
			`INSERT INTO gauges (id, g_value)
			VALUES ($1, $2)
			ON CONFLICT (id)
			DO UPDATE SET g_value = $2
			RETURNING id, g_value;
			`,
			m.ID,
			*m.Value,
		)
		err := row.Scan(&id, &gValue)
		if err != nil {
			if checkRetry(err) {
				err = storagerror.NewRetriableError(err)
			}
			return metric.Metric{}, fmt.Errorf("failed to scan row when update metric: %w", err)
		}
		updating = metric.New(metric.TypeGauge)
		updating.ID = id
		updating.Value = &gValue
		return updating, nil
	}

	var cValue int64
	row = psg.pool.QueryRow(
		context.TODO(),
		`INSERT INTO counters (id, delta)
			VALUES ($1, $2)
			ON CONFLICT (id)
			DO UPDATE SET delta = counters.delta + EXCLUDED.delta
			RETURNING id, delta;`,
		m.ID,
		*m.Delta,
	)
	err := row.Scan(&id, &cValue)
	if err != nil {
		if checkRetry(err) {
			err = storagerror.NewRetriableError(err)
		}
		return metric.Metric{}, fmt.Errorf("failed to scan row when update metric: %w", err)
	}
	updating = metric.New(metric.TypeCounter)
	updating.ID = id
	updating.Delta = &cValue

	return updating, nil
}

func (psg *PostgresStorage) UpdateManyMetrics(ctx context.Context, mList []metric.Metric) error {
	log := psg.log.LogrusLog
	tx, err := psg.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed start a transaction: %w", err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			if !errors.Is(err, pgx.ErrTxClosed) {
				log.Errorf("failed rolls back the transaction: %v", err)
			}
		}
	}()

	log.Info("updating metrics ...")

	for _, m := range mList {
		switch m.MType {
		case metric.TypeCounter:
			log.WithFields(logrus.Fields{
				"ID":    m.ID,
				"Type":  m.MType,
				"Delta": *m.Delta,
			}).Info("metric for updating")
			_, err = tx.Exec(ctx,
				`INSERT INTO counters (id, delta) 
				VALUES($1, $2)
				ON CONFLICT (id)
				DO UPDATE SET delta = counters.delta + $2;`,
				m.ID,
				*m.Delta,
			)
		case metric.TypeGauge:
			log.WithFields(logrus.Fields{
				"ID":    m.ID,
				"Type":  m.MType,
				"Value": *m.Value,
			}).Info("metric for updating")
			_, err = tx.Exec(ctx,
				`INSERT INTO gauges (id, g_value) 
				VALUES($1, $2)
				ON CONFLICT (id)
				DO UPDATE SET g_value = $2;`,
				m.ID,
				*m.Value,
			)
		default:
			return errors.New("failed update metrics: unknown metric type")
		}
		if err != nil {
			if checkRetry(err) {
				err = storagerror.NewRetriableError(err)
			}
			return fmt.Errorf("failed exec query update metric: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed commits the transaction update metrics: %w", err)
	}
	return nil
}

func (psg *PostgresStorage) GetAllMetrics() ([][2]string, error) {
	res := make([][2]string, 0)

	rowsGauge, err := psg.pool.Query(
		context.TODO(),
		`SELECT id, g_value FROM gauges;`,
	)
	if err != nil {
		if checkRetry(err) {
			err = storagerror.NewRetriableError(err)
		}
		return res, fmt.Errorf("failed to select all metrics from gauges table: %w", err)
	}
	defer rowsGauge.Close()

	var id string
	var gVal float64
	for rowsGauge.Next() {
		if err := rowsGauge.Scan(&id, &gVal); err != nil {
			if checkRetry(err) {
				err = storagerror.NewRetriableError(err)
			}
			return res, fmt.Errorf("failed to scan gauge metric when get all metrics: %w", err)
		}
		res = append(res, [2]string{id, strconv.FormatFloat(gVal, 'g', -1, 64)})
	}

	rowsCounter, err := psg.pool.Query(
		context.TODO(),
		`SELECT id, delta FROM counters;`,
	)
	if err != nil {
		if checkRetry(err) {
			err = storagerror.NewRetriableError(err)
		}
		return res, fmt.Errorf("failed to select all metrics from counters table: %w", err)
	}
	defer rowsCounter.Close()

	var cVal int64
	for rowsCounter.Next() {
		if err := rowsCounter.Scan(&id, &cVal); err != nil {
			if checkRetry(err) {
				err = storagerror.NewRetriableError(err)
			}
			return res, fmt.Errorf("failed to scan counter metric when get all metrics: %w", err)
		}
		res = append(res, [2]string{id, strconv.FormatInt(cVal, 10)})
	}

	return res, nil
}

func (psg *PostgresStorage) GetMetricValue(name, typeMetric string) (metric.Metric, error) {
	var resMetric metric.Metric
	var row pgx.Row
	var id string
	if typeMetric == metric.TypeGauge {
		var gValue float64
		row = psg.pool.QueryRow(
			context.TODO(),
			`SELECT id, g_value FROM gauges
			WHERE id = $1;
			`,
			name,
		)
		err := row.Scan(&id, &gValue)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return metric.Metric{}, metric.ErrUnknownMetric
			}
			if checkRetry(err) {
				err = storagerror.NewRetriableError(err)
			}
			return metric.Metric{}, fmt.Errorf("failed to scan row when get metric: %w", err)
		}
		resMetric = metric.New(metric.TypeGauge)
		resMetric.ID = id
		resMetric.Value = &gValue
		return resMetric, nil
	}
	var cValue int64
	row = psg.pool.QueryRow(
		context.TODO(),
		`SELECT id, delta FROM counters
			WHERE id = $1`,
		name,
	)
	err := row.Scan(&id, &cValue)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return metric.Metric{}, metric.ErrUnknownMetric
		}
		if checkRetry(err) {
			err = storagerror.NewRetriableError(err)
		}
		return metric.Metric{}, fmt.Errorf("failed to scan row when get metric: %w", err)
	}
	resMetric = metric.New(metric.TypeCounter)
	resMetric.ID = id
	resMetric.Delta = &cValue
	return resMetric, nil
}

func (psg *PostgresStorage) Close() error {
	psg.pool.Close()
	return nil
}

func checkRetry(err error) bool {
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
