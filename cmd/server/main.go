package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/zhenyanesterkova/metricsmonitor/handlers"
	"github.com/zhenyanesterkova/metricsmonitor/storage/memstorage"
)

func NewRouter(storage handlers.Repositorie) chi.Router {
	router := chi.NewRouter()

	router.Route("/", func(r chi.Router) {
		r.Get("/", handlers.New("getAllMetrics", storage))
		r.Route("/value", func(r chi.Router) {
			r.Get("/{typeMetric}/{nameMetric}", handlers.New("getMetricValue", storage))
		})
		r.Route("/update", func(r chi.Router) {
			r.Post("/{typeMetric}/{nameMetric}/{valueMetric}", handlers.New("updateMetricValue", storage))
		})

	})

	return router
}

func main() {

	parseFlags()

	storage := memstorage.New()

	if err := http.ListenAndServe(endpoint, NewRouter(storage)); err != nil {
		panic(err)
	}

}
