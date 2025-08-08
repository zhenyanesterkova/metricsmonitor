package server

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/backoff"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/retrystorage"
)

type BuildInfo struct {
	Version string
	Date    string
	Commit  string
}

type ServerType int

const (
	HTTPServer ServerType = iota
	GRPCServer
)

type Server struct {
	config       *config.Config
	logger       logger.LogrusLogger
	storage      *retrystorage.RetryStorage
	shutdownChan chan os.Signal
	buildInfo    BuildInfo
	serverType   ServerType
}

func New(serverType ServerType, buildInfo BuildInfo) *Server {
	return &Server{
		serverType:   serverType,
		buildInfo:    buildInfo,
		shutdownChan: make(chan os.Signal, 1),
	}
}

func (s *Server) Initialize() error {
	s.config = config.New()
	if err := s.config.Build(); err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	s.logger = logger.NewLogrusLogger()
	if err := s.logger.SetLevelForLog(s.config.LConfig.Level); err != nil {
		s.logger.LogrusLog.Errorf("can not parse log level: %v", err)
		return fmt.Errorf("parse log level error: %w", err)
	}

	if err := s.initStorage(); err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	s.logBuildInfo()

	return nil
}

func (s *Server) initStorage() error {
	backoffInst := backoff.New(
		s.config.RetryConfig.MinDelay,
		s.config.RetryConfig.MaxDelay,
		s.config.RetryConfig.MaxAttempt,
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

	retryStore, err := retrystorage.New(s.config.DBConfig, s.logger, backoffInst, checkRetryFunc)
	if err != nil {
		s.logger.LogrusLog.Errorf("failed create storage: %v", err)
		return fmt.Errorf("failed create storage: %w", err)
	}

	s.storage = retryStore
	return nil
}

func (s *Server) logBuildInfo() {
	s.logger.LogrusLog.Infof("Build version: %s", s.buildInfo.Version)
	s.logger.LogrusLog.Infof("Build date: %s", s.buildInfo.Date)
	s.logger.LogrusLog.Infof("Build commit: %s", s.buildInfo.Commit)
}

func (s *Server) Run() error {
	switch s.serverType {
	case HTTPServer:
		return s.runHTTP()
	case GRPCServer:
		return s.runGRPC()
	default:
		return fmt.Errorf("unknown server type: %d", s.serverType)
	}
}

func (s *Server) waitForShutdown(cleanup func() error) error {
	signal.Notify(s.shutdownChan, os.Interrupt, syscall.SIGTERM)

	sig := <-s.shutdownChan
	s.logger.LogrusLog.Info("Got signal: ", sig)

	if cleanup != nil {
		if err := cleanup(); err != nil {
			s.logger.LogrusLog.Errorf("error during cleanup: %v", err)
			return fmt.Errorf("error during cleanup: %w", err)
		}
	}

	return nil
}

func (s *Server) Close() error {
	if s.storage != nil {
		if err := s.storage.Close(); err != nil {
			s.logger.LogrusLog.Errorf("can not close storage: %v", err)
			return fmt.Errorf("close storage error: %w", err)
		}
	}
	return nil
}
