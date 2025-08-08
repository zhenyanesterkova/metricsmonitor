package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zhenyanesterkova/metricsmonitor/internal/handler"
)

func (s *Server) runHTTP() error {
	router := chi.NewRouter()

	repoHandler := handler.NewRepositorieHandler(
		s.storage,
		s.logger,
		s.config.SConfig.HashKey,
		s.config.SConfig.CryptoPrivateKeyPath,
		s.config.SConfig.TIpNet,
	)

	if err := repoHandler.InitChiRouter(router); err != nil {
		s.logger.LogrusLog.Errorf("can not init router: %v", err)
		return fmt.Errorf("failed init router: %w", err)
	}

	s.logger.LogrusLog.Infof("Start HTTP Server on %s", s.config.SConfig.Address)

	server := &http.Server{
		Addr:    s.config.SConfig.Address,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.LogrusLog.Errorf("HTTP server error: %v", err)
		}
	}()

	return s.waitForShutdown(func() error {
		return server.Shutdown(context.Background())
	})
}
