package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/middlewares"
)

type Repositorie interface {
	UpdateMetric(name string, typeMetric string, val string) error
	GetAllMetrics() ([][2]string, error)
	GetMetricValue(name, metricType string) (string, error)
}

type RepositorieHandler struct {
	Repo   Repositorie
	Logger *logrus.Logger
}

func NewRepositorieHandler(rep Repositorie) *RepositorieHandler {
	return &RepositorieHandler{
		Repo:   rep,
		Logger: logger.Logger(),
	}

}

func (rh *RepositorieHandler) InitChiRouter(router *chi.Mux) {
	router.Use(middlewares.RequestLogger)
	router.Route("/", func(r chi.Router) {

		r.Get("/", rh.GetAllMetrics)

		r.Get("/value/{typeMetric}/{nameMetric}", rh.GetMetricValue)

		r.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", rh.UpdateMetric)

	})
}
