package main

import (
	"net/http"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/handlers"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memstorage"

	"github.com/go-chi/chi/v5"
)

func main() {

	cfgBuilder := config.GetConfigBuilder()
	cfg := cfgBuilder.Build()

	storage := memstorage.New()

	router := chi.NewRouter()

	repoHandler := handlers.NewRepositorieHandler(storage)
	repoHandler.InitChiRouter(router)

	if err := http.ListenAndServe(cfg.SConfig.Address, router); err != nil {
		panic(err)
	}

}
