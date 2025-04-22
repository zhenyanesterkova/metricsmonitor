// Package handler contains a set of HTTP handlers for processing incoming requests.
//
// Package structure:
//
// - The main handlers are organized as methods of the RepositorieHandler structure.
//
// - Dependencies are injected through the constructor function NewRepositorieHandler().
//
// - The functionality can be extended through middleware.
package handler

import (
	"context"
	"fmt"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/middleware"
)

// TextServerError - represents the error text,
// which is sent to the user when
// unexpected internal server errors occur.
const (
	TextServerError = "Something went wrong... Server error"
)

// Repositorie represents an interface for working with data storage systems.
// It provides unified access to various types of storage,
// including databases, file systems, and other storage systems.
type Repositorie interface {
	// UpdateMetric updates a metric with the given type and value.
	UpdateMetric(metric.Metric) (metric.Metric, error)
	// GetAllMetrics retrieves all available metrics from the storage.
	GetAllMetrics() ([][2]string, error)
	// GetMetricValue retrieves a specific metric from the storage by its name and type.
	GetMetricValue(name, typeMetric string) (metric.Metric, error)
	// Ping checks the availability of the storage.
	Ping() error
	// UpdateManyMetrics updates multiple metrics in the storage.
	UpdateManyMetrics(ctx context.Context, mList []metric.Metric) error
}

// RepositorieHandler provides methods to handle various operations.
type RepositorieHandler struct {
	// Repo - data storage.
	Repo Repositorie
	// Logger is a logging utility used to record events and errors.
	Logger logger.LogrusLogger
	// hashKey is a key for calculating the hash.
	hashKey *string
	// pathToPrivateKey used to store the path to a file containing an asymmetric encryption private key
	pathToPrivateKey string
}

// NewRepositorieHandler creates a new instance of RepositoryHandler.
func NewRepositorieHandler(
	rep Repositorie,
	log logger.LogrusLogger,
	key *string,
	pathToPrivateKey string,
) *RepositorieHandler {
	return &RepositorieHandler{
		Repo:             rep,
		Logger:           log,
		hashKey:          key,
		pathToPrivateKey: pathToPrivateKey,
	}
}

// InitChiRouter initializes a new Chi router with predefined routes and middleware.
func (rh *RepositorieHandler) InitChiRouter(router *chi.Mux) error {
	mdlWare, err := middleware.NewMiddlewareStruct(rh.Logger, rh.hashKey, rh.pathToPrivateKey)
	if err != nil {
		return fmt.Errorf("failed create struct for middleware: %w", err)
	}
	router.Use(mdlWare.ResetRespDataStruct)
	router.Use(mdlWare.RequestLogger)
	if rh.hashKey != nil {
		router.Use(mdlWare.CheckSignData)
	}
	router.Use(mdlWare.DecryptionMiddleware)
	router.Use(mdlWare.GZipMiddleware)
	router.Route("/", func(r chi.Router) {
		r.Get("/", rh.GetAllMetrics)
		r.Get("/ping", rh.Ping)
		r.Route("/value/", func(r chi.Router) {
			r.Post("/", rh.GetMetricValueJSON)
			r.Get("/{typeMetric}/{nameMetric}", rh.GetMetricValue)
		})

		r.Route("/updates/", func(r chi.Router) {
			r.Post("/", rh.UpdateManyMetrics)
		})
		r.Route("/update/", func(r chi.Router) {
			r.Post("/", rh.UpdateMetricJSON)
			r.Post("/{typeMetric}/{nameMetric}/{valueMetric}", rh.UpdateMetric)
		})
		r.Route("/debug/pprof/", func(r chi.Router) {
			r.Get("/", pprof.Index)
			r.Get("/cmdline", pprof.Cmdline)
			r.Get("/symbol", pprof.Symbol)
			r.Get("/trace", pprof.Trace)
			r.Get("/profile", pprof.Profile)
			r.Handle("/goroutine", pprof.Handler("goroutine"))
			r.Handle("/threadcreate", pprof.Handler("threadcreate"))
			r.Handle("/heap", pprof.Handler("heap"))
			r.Handle("/block", pprof.Handler("block"))
			r.Handle("/mutex", pprof.Handler("mutex"))
			r.Handle("/allocs", pprof.Handler("allocs"))
		})
	})
	return nil
}
