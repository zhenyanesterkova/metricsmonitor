package handler

import (
	"context"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/middleware"
)

const (
	TextServerError = "Something went wrong... Server error"
)

type Repositorie interface {
	UpdateMetric(metric.Metric) (metric.Metric, error)
	GetAllMetrics() ([][2]string, error)
	GetMetricValue(name, typeMetric string) (metric.Metric, error)
	Ping() error
	UpdateManyMetrics(ctx context.Context, mList []metric.Metric) error
}

type RepositorieHandler struct {
	Repo    Repositorie
	Logger  logger.LogrusLogger
	hashKey *string
}

func NewRepositorieHandler(
	rep Repositorie,
	log logger.LogrusLogger,
	key *string,
) *RepositorieHandler {
	return &RepositorieHandler{
		Repo:    rep,
		Logger:  log,
		hashKey: key,
	}
}

func (rh *RepositorieHandler) InitChiRouter(router *chi.Mux) {
	mdlWare := middleware.NewMiddlewareStruct(rh.Logger, rh.hashKey)
	router.Use(mdlWare.ResetRespDataStruct)
	router.Use(mdlWare.RequestLogger)
	if rh.hashKey != nil {
		router.Use(mdlWare.CheckSignData)
	}
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
		r.Route("/debug/pprof/", func(r chi.Router) {
			r.Get("/", pprof.Index)
			r.Get("/cmdline", pprof.Cmdline)
			r.Get("/symbol", pprof.Symbol)
			r.Get("/trace", pprof.Trace)
			r.Route("/profile", func(r chi.Router) {
				r.Get("/", pprof.Profile)
			})
			r.Handle("/goroutine", pprof.Handler("goroutine"))
			r.Handle("/threadcreate", pprof.Handler("threadcreate"))
			r.Handle("/heap", pprof.Handler("heap"))
			r.Handle("/block", pprof.Handler("block"))
			r.Handle("/mutex", pprof.Handler("mutex"))
			r.Handle("/allocs", pprof.Handler("allocs"))
		})
	})
}
