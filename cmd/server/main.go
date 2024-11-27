package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/rwfile"
	"github.com/zhenyanesterkova/metricsmonitor/internal/handler"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memstorage"
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

	storage := memstorage.New()

	fileWriter, err := rwfile.NewFileWriter(cfg.RConfig.FileStoragePath)
	if err != nil {
		loggerInst.LogrusLog.Errorf("can not create file writer: %v", err)
		return fmt.Errorf("file writer error: %w", err)
	}
	fileReader, err := rwfile.NewFileReader(cfg.RConfig.FileStoragePath)
	if err != nil {
		loggerInst.LogrusLog.Errorf("can not create file reader: %v", err)
		return fmt.Errorf("file reader error: %w", err)
	}
	defer func() {
		err := fileWriter.Close()
		if err != nil {
			loggerInst.LogrusLog.Errorf("can not close file writer: %v", err)
		}
	}()
	defer func() {
		err := fileReader.Close()
		if err != nil {
			loggerInst.LogrusLog.Errorf("can not close file reader: %v", err)
		}
	}()
	if cfg.RConfig.Restore {
		mementoMemStorage := storage.CreateMemento()
		err := fileReader.ReadSnapStorage(mementoMemStorage)
		if err != nil {
			loggerInst.LogrusLog.Errorf("can not read snapshot of storage from file %s: %v", cfg.RConfig.FileStoragePath, err)
			return fmt.Errorf("snapshot of storage error: %w", err)
		}
		storage.RestoreMemento(mementoMemStorage)
	}

	router := chi.NewRouter()

	repoHandler := handler.NewRepositorieHandler(storage, loggerInst)
	repoHandler.InitChiRouter(router)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	loggerInst.LogrusLog.Debugf("Start Server on %s", cfg.SConfig.Address)
	loggerInst.LogrusLog.Infof("Start Server on %s", cfg.SConfig.Address)
	go func() {
		if err := http.ListenAndServe(cfg.SConfig.Address, router); err != nil {
			loggerInst.LogrusLog.Errorf("server error: %v", err)
		}
	}()
	go func() {
		ticker := time.NewTicker(cfg.RConfig.StoreInterval)
		for range ticker.C {
			loggerInst.LogrusLog.Info("starting storage copying...")
			mementoStorage := storage.CreateMemento()
			err := fileWriter.WriteSnapStorage(*mementoStorage)
			if err != nil {
				loggerInst.LogrusLog.Errorf("can not write snapshot of storage to file %s: %v", cfg.RConfig.FileStoragePath, err)
			}
			loggerInst.LogrusLog.Info("end storage copying...")
		}
	}()

	s := <-c
	loggerInst.LogrusLog.Info("Got signal: ", s)

	loggerInst.LogrusLog.Info("starting storage copying...")
	mementoStorage := storage.CreateMemento()
	err = fileWriter.WriteSnapStorage(*mementoStorage)
	if err != nil {
		loggerInst.LogrusLog.Errorf("can not write snapshot of storage to file %s: %v", cfg.RConfig.FileStoragePath, err)
		return fmt.Errorf("snapshot of storage error when exit from app: %w", err)
	}
	loggerInst.LogrusLog.Info("end storage copying...")

	return nil
}
