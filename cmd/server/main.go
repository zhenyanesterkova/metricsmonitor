package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/backoff"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/handler"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func run() error {
	cfg := config.New()
	err := cfg.Build()
	if err != nil {
		log.Printf("can not build config: %v", err)
		return fmt.Errorf("config error: %w", err)
	}

	loggerInst := logger.NewLogrusLogger()
	err = loggerInst.SetLevelForLog(cfg.LConfig.Level)
	if err != nil {
		loggerInst.LogrusLog.Errorf("can not parse log level: %v", err)
		return fmt.Errorf("parse log level error: %w", err)
	}

	store, err := storage.NewStore(cfg.DBConfig, loggerInst)
	if err != nil {
		loggerInst.LogrusLog.Errorf("can not create storage: %v", err)
		return fmt.Errorf("can not create storage: %w", err)
	}

	defer func() {
		err := store.Close()
		if err != nil {
			loggerInst.LogrusLog.Errorf("can not close storage: %v", err)
		}
	}()

	backoffInst := backoff.New(
		cfg.RetryConfig.MinDelay,
		cfg.RetryConfig.MaxDelay,
		cfg.RetryConfig.MaxAttempt,
	)

	router := chi.NewRouter()

	repoHandler := handler.NewRepositorieHandler(store, loggerInst, backoffInst)
	repoHandler.InitChiRouter(router)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	loggerInst.LogrusLog.Infof("Start Server on %s", cfg.SConfig.Address)
	go func() {
		if err := http.ListenAndServe(cfg.SConfig.Address, router); err != nil {
			loggerInst.LogrusLog.Errorf("server error: %v", err)
		}
	}()

	s := <-c
	loggerInst.LogrusLog.Info("Got signal: ", s)

	return nil
}
