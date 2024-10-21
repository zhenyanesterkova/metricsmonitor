package handler

import (
	"github.com/go-chi/chi/v5"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/middleware"
)

type Repositorie interface {
	UpdateMetric(name string, typeMetric string, val string) error
	GetAllMetrics() ([][2]string, error)
	GetMetricValue(name, metricType string) (string, error)
}

type RepositorieHandler struct {
	Repo   Repositorie
	Logger logger.LogrusLogger
}

func NewRepositorieHandler(rep Repositorie, log logger.LogrusLogger) *RepositorieHandler {
	return &RepositorieHandler{
		Repo:   rep,
		Logger: log,
	}

}

func (rh *RepositorieHandler) InitChiRouter(router *chi.Mux) {
	mdlWare := middleware.NewLoggerMiddleware(rh.Logger)
	router.Use(mdlWare.RequestLogger)
	router.Route("/", func(r chi.Router) {

		r.Get("/", rh.GetAllMetrics)

		r.Get("/value/{typeMetric}/{nameMetric}", rh.GetMetricValue)

		r.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", rh.UpdateMetric)

	})
}
