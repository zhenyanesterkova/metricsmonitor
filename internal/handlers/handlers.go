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

func NewRepositorieHandler(rep Repositorie) *RepositorieHandler {
	return &RepositorieHandler{
		Repo: rep,
	}

}

func (rh *RepositorieHandler) InitChiRouter(router *chi.Mux) {
	router.Route("/", func(r chi.Router) {

		r.Get("/", rh.GetAllMetrics())

		r.Get("/value/{typeMetric}/{nameMetric}", rh.GetMetricValue())

		r.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", rh.UpdateMetric())

	})
}
