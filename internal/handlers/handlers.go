package handlers

import (
	"net/http"
)

type Repositorie interface {
	UpdateMetric(name string, typeMetric string, val string) error
	GetAllMetrics() ([][2]string, error)
	GetMetricValue(name, metricType string) (string, error)
}

func New(handlerName string, rep Repositorie) http.HandlerFunc {
	switch handlerName {
	case "updateMetricValue":
		return updateMetricValue(rep)
	case "getMetricValue":
		return getMetricValue(rep)
	case "getAllMetrics":
		return getAllMetrics(rep)
	default:
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "unknown handler name", http.StatusInternalServerError)
		}
	}
}
