package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/handler"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memstorage"
)

func main() {

	cfgBuilder := config.GetConfigBuilder()
	cfg, err := cfgBuilder.Build()
	if err != nil {
		fmt.Printf("can not build config: %v", err)
		os.Exit(1)
	}

	logger := logger.NewLogrusLogger()
	err = logger.SetLevelForLog(cfg.LConfig.Level)
	if err != nil {
		logger.LogrusLog.Errorf("can not parse log level: %v", err)
	}

	storage := memstorage.New()

	router := chi.NewRouter()

	repoHandler := handler.NewRepositorieHandler(storage, logger)
	repoHandler.InitChiRouter(router)

	logger.LogrusLog.Debugf("Start Server on %s", cfg.SConfig.Address)
	if err := http.ListenAndServe(cfg.SConfig.Address, router); err != nil {
		panic(err)
	}

}
