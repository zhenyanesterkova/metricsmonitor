package handlers

import (
	"github.com/go-chi/chi/v5"
)

type Repositorie interface {
	UpdateMetric(name string, typeMetric string, val string) error
	GetAllMetrics() ([][2]string, error)
	GetMetricValue(name, metricType string) (string, error)
}

type RepositorieHandler struct {
	Repo Repositorie
}

func NewRepositorieHandler(router *chi.Mux, rep Repositorie) {
	handler := &RepositorieHandler{
		Repo: rep,
	}

	router.Route("/", func(r chi.Router) {

		r.Get("/", handler.GetAllMetrics())

		r.Get("/value/{typeMetric}/{nameMetric}", handler.GetMetricValue())

		r.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", handler.UpdateMetric())

	})
}
