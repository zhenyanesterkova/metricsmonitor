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

func main() {

	cfg := getConfig()

	storage := memstorage.New()

	router := chi.NewRouter()

	repoHandler := handlers.NewRepositorieHandler(storage)
	repoHandler.InitChiRouter(router)

	if err := http.ListenAndServe(cfg.SConfig.Address, router); err != nil {
		panic(err)
	}

}
