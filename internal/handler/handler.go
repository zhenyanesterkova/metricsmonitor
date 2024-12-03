package handler

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/retry"
	"github.com/zhenyanesterkova/metricsmonitor/internal/middleware"
)

const (
	TextServerError = "Something went wrong... Server error"
)

type Repositorie interface {
	UpdateMetric(metric.Metric) (metric.Metric, error)
	GetAllMetrics() ([][2]string, error)
	GetMetricValue(name, typeMetric string) (metric.Metric, error)
	Ping() (bool, error)
	UpdateManyMetrics(ctx context.Context, mList []metric.Metric) error
}

type RepositorieHandler struct {
	Repo    Repositorie
	retrier retry.Retrier
	Logger  logger.LogrusLogger
	DSN     string
}

func NewRepositorieHandler(
	rep Repositorie,
	log logger.LogrusLogger,
	dsn string, retCfg config.RetryConfig,
) *RepositorieHandler {
	b := retry.NewBackoff(retCfg.Min, retCfg.Max, retCfg.MaxAttempt)
	return &RepositorieHandler{
		Repo:    rep,
		Logger:  log,
		DSN:     dsn,
		retrier: retry.NewRetrier(b, nil),
	}
}

func (rh *RepositorieHandler) InitChiRouter(router *chi.Mux) {
	mdlWare := middleware.NewMiddlewareStruct(rh.Logger)
	router.Use(mdlWare.RequestLogger)
	router.Use(mdlWare.GZipMiddleware)
	router.Route("/", func(r chi.Router) {
		r.Get("/", rh.GetAllMetrics)
		r.Get("/ping", rh.Ping)
		r.Route("/value/", func(r chi.Router) {
			r.Post("/", rh.GetMetricValueJSON)
			r.Get("/{typeMetric}/{nameMetric}", rh.GetMetricValue)
		})

		r.Route("/updates/", func(r chi.Router) {
			r.Post("/", rh.UpdateManyMetrics)
		})
		r.Route("/update/", func(r chi.Router) {
			r.Post("/", rh.UpdateMetricJSON)
			r.Post("/{typeMetric}/{nameMetric}/{valueMetric}", rh.UpdateMetric)
		})
	})
}
