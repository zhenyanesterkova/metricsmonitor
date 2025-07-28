package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	proto "github.com/zhenyanesterkova/metricsmonitor/internal/app/proto/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/backoff"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/grpchandler"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/retrystorage"
	"google.golang.org/grpc"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {
	if err := run(); err != nil {
		log.Fatalf("grpc server error: %v", err)
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

	backoffInst := backoff.New(
		cfg.RetryConfig.MinDelay,
		cfg.RetryConfig.MaxDelay,
		cfg.RetryConfig.MaxAttempt,
	)

	checkRetryFunc := func(err error) bool {
		var pgErr *pgconn.PgError
		var pgErrConn *pgconn.ConnectError
		res := false
		if errors.As(err, &pgErr) {
			res = pgerrcode.IsConnectionException(pgErr.Code)
		} else if errors.As(err, &pgErrConn) {
			res = true
		}
		return res
	}

	retryStore, err := retrystorage.New(cfg.DBConfig, loggerInst, backoffInst, checkRetryFunc)
	if err != nil {
		loggerInst.LogrusLog.Errorf("failed create storage: %v", err)
		return fmt.Errorf("failed create storage: %w", err)
	}
	defer func() {
		err := retryStore.Close()
		if err != nil {
			loggerInst.LogrusLog.Errorf("can not close storage: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	loggerInst.LogrusLog.Infof("Build version: %s\n", buildVersion)
	loggerInst.LogrusLog.Infof("Build date: %s\n", buildDate)
	loggerInst.LogrusLog.Infof("Build commit: %s\n", buildCommit)

	listen, err := net.Listen("tcp", cfg.SConfig.Address)
	if err != nil {
		loggerInst.LogrusLog.Errorf("failed listen announces on the local network address: %v", err)
	}

	s := grpc.NewServer()

	grpcHandler := grpchandler.New(
		retryStore,
		loggerInst,
	)
	proto.RegisterMonitorServer(s, grpcHandler)

	loggerInst.LogrusLog.Infof("Start gRPC Server on %s", cfg.SConfig.Address)
	go func() {
		if err := s.Serve(listen); err != nil {
			loggerInst.LogrusLog.Errorf("grpc server error: %v", err)
		}
	}()

	sig := <-c
	loggerInst.LogrusLog.Info("Got signal: ", sig)

	return nil
}
