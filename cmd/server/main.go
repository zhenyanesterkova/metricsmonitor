package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zhenyanesterkova/metricsmonitor/internal/handlers"
	"github.com/zhenyanesterkova/metricsmonitor/internal/handlers/storage/get"
	"github.com/zhenyanesterkova/metricsmonitor/internal/handlers/storage/update"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memstorage"
)

var storage *memstorage.Storage

func NewRouter(storage handlers.Storage) chi.Router {
	router := chi.NewRouter()

	router.Route("/", func(r chi.Router) {
		r.Get("/", get.NewGetAll(storage))
		r.Get("/value/{typeMetric}/{nameMetric}", get.NewGet(storage))
		r.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", update.New(storage))

	})

	return router
}

func main() {
	parseFlags()

	storage = memstorage.New()

	if err := http.ListenAndServe(endpoint, NewRouter(storage)); err != nil {
		panic(err)
	}

}
