package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/handlers"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memstorage"
)

func main() {

	cfgBuilder := config.GetConfigBuilder()
	cfg := cfgBuilder.Build()

	log := logger.Logger()
	logger.SetLevelForLog(cfg.LConfig.Level)

	storage := memstorage.New()

	router := chi.NewRouter()

	repoHandler := handlers.NewRepositorieHandler(storage)
	repoHandler.InitChiRouter(router)

	log.Debugf("Start Server on %s", cfg.SConfig.Address)
	if err := http.ListenAndServe(cfg.SConfig.Address, router); err != nil {
		panic(err)
	}

}
