package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/handlers"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memstorage"
)

func getConfig() config.Config {
	cfgBuilder := config.GetConfigBuilder()
	cfgDirector := config.NewConfigDirector(cfgBuilder)
	resConfig := cfgDirector.BuildConfig()

	return resConfig
}

func NewRouter(storage handlers.Repositorie) chi.Router {
	router := chi.NewRouter()

	router.Route("/", func(r chi.Router) {
		r.Get("/", handlers.New("getAllMetrics", storage))

		r.Get("/value/{typeMetric}/{nameMetric}", handlers.New("getMetricValue", storage))

		r.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", handlers.New("updateMetricValue", storage))

	})

	return router
}

func main() {

	cfg := getConfig()

	storage := memstorage.New()

	if err := http.ListenAndServe(cfg.SConfig.Address, NewRouter(storage)); err != nil {
		panic(err)
	}

}
