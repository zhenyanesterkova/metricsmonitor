package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/handler"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memstorage"
)

func main() {

	cfgBuilder := config.GetConfigBuilder()
	cfg := cfgBuilder.Build()

	log := logger.Logger()
	err := logger.SetLevelForLog(cfg.LConfig.Level)
	if err != nil {
		log.Errorf("can not parse log level: %v", err)
	}

	storage := memstorage.New()

	router := chi.NewRouter()

	repoHandler := handler.NewRepositorieHandler(storage)
	repoHandler.InitChiRouter(router)

	log.Debugf("Start Server on %s", cfg.SConfig.Address)
	if err := http.ListenAndServe(cfg.SConfig.Address, router); err != nil {
		panic(err)
	}

}
