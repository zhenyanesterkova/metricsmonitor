package handler

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
)

// UpdateMetric handles HTTP requests to update a specific metric value.
//
// URL Parameters:
//   - typeMetric: The type of the metric (gauge, counter).
//   - nameMetric: The name of the metric to update.
//   - valueMetric: The new value for the metric.
//
// Supported Metric Types:
//   - Gauge: Floating value.
//   - Counter: Integer value.
//
// HTTP Response:
//
// Sets appropriate HTTP status codes:
//
//   - 200 OK: Successful update.
//   - 400 Bad Request: Invalid metric type or value.
//   - 404 Not Found: Invalid metric name.
//   - 500 Internal Server Error: Other errors.
//
// Example usage in HTTP request:
//
//	POST /update/gauge/Alloc/22.5
//	POST /update/counter/PollCount/3
func (rh *RepositorieHandler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	log := rh.Logger.LogrusLog

	log.Info("updating metric ...")

	metricType := chi.URLParam(r, "typeMetric")
	metricName := chi.URLParam(r, "nameMetric")
	metricValue := chi.URLParam(r, "valueMetric")

	metrica := metric.New(metricType)
	metrica.ID = metricName
	switch metricType {
	case metric.TypeGauge:
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		*metrica.Value = val
		log.WithFields(logrus.Fields{
			"ID":    metrica.ID,
			"Type":  metrica.MType,
			"Value": *metrica.Value,
		}).Info("gauge metric updating")
	case metric.TypeCounter:
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		*metrica.Delta = val
		log.WithFields(logrus.Fields{
			"ID":    metrica.ID,
			"Type":  metrica.MType,
			"Delta": *metrica.Delta,
		}).Info("counter metric updating")
	}

	_, err := rh.Repo.UpdateMetric(metrica)
	if err != nil {
		switch {
		case errors.Is(err, metric.ErrInvalidName):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, metric.ErrParseValue) ||
			errors.Is(err, metric.ErrUnknownType) ||
			errors.Is(err, metric.ErrInvalidType):
			w.WriteHeader(http.StatusBadRequest)
		default:
			log.Errorf("handler func UpdateMetric(): error update metric - %v", err)
			http.Error(w, TextServerError, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateMetricJSON handles HTTP requests to update a specific metric value via JSON payload.
//
// HTTP Request:
// - Accepts a JSON payload with metric details.
// - Example request body:
//
// for counter metric:
//
//	{
//	  "ID": "metric_name",
//	  "MType": "metric_type",
//	  "Delta": 100
//	}
//
// for gauge metric:
//
//	{
//	  "ID": "metric_name",
//	  "MType": "metric_type",
//	  "Value": 5.5
//	}
//
// HTTP Response:
//
// - Returns the updated metric in JSON format.
//
// - Sets appropriate HTTP status codes:
//   - 200 OK: Successful update
//   - 400 Bad Request: Invalid metric type or value
//   - 404 Not Found: Invalid metric name
//   - 500 Internal Server Error: Other errors
//
// Example usage in HTTP request:
//
//	POST /update/
//	Content-Type: application/json
//
//	{
//	  "ID": "TotalAlloc",
//	  "MType": "gauge",
//	  "Value": 22.5
//	}
//
// Example response:
//
//	{
//	  "ID": "TotalAlloc",
//	  "MType": "gauge",
//	  "Value": 22.5
//	}
func (rh *RepositorieHandler) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	log := rh.Logger.LogrusLog

	if rh.trustIPNet != nil {
		ipStr := r.Header.Get("X-Real-IP")

		ip := net.ParseIP(ipStr)

		allowed := rh.trustIPNet.Contains(ip)

		if !allowed {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	log.Info("updating metric ...")

	newMetric := metric.New("")
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&newMetric); err != nil {
		log.Errorf("handler func UpdateMetricJSON(): error decode metric - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}

	log.WithFields(logrus.Fields{
		"ID":    newMetric.ID,
		"Type":  newMetric.MType,
		"Value": *newMetric.Value,
		"Delta": *newMetric.Delta,
	}).Info("metric for updating")

	updating, err := rh.Repo.UpdateMetric(newMetric)
	if err != nil {
		switch {
		case errors.Is(err, metric.ErrInvalidName):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, metric.ErrParseValue) ||
			errors.Is(err, metric.ErrUnknownType) ||
			errors.Is(err, metric.ErrInvalidType):
			w.WriteHeader(http.StatusBadRequest)
		default:
			log.Errorf("handler func UpdateMetricJSON(): error update metric - %v", err)
			http.Error(w, TextServerError, http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(updating); err != nil {
		log.Errorf("handler func UpdateMetricJSON(): error encode metric - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}
}

// UpdateManyMetrics handles HTTP requests to update multiple metrics in a single request.
//
// HTTP Request:
//
// - Accepts a JSON array of metric objects
//
// - Example request body:
//
//	[
//	    {
//	        "ID": "metric1",
//	        "MType": "gauge",
//	        "Value": 22.5
//	    },
//	    {
//	        "ID": "metric2",
//	        "MType": "counter",
//	        "Delta": 100
//	    }
//	]
//
// HTTP Response:
//   - Returns 200 OK if all metrics are updated successfully
//   - Returns 500 Internal Server Error if there's an error during processing
//
// Example usage in HTTP request:
//
//	POST /updates/
//	Content-Type: application/json
//
//	[
//	    {
//	        "ID": "TotalAlloc",
//	        "MType": "gauge",
//	        "Value": 22.5
//	    },
//	    {
//	        "ID": "PollCount",
//	        "MType": "counter",
//	        "Delta": 100
//	    }
//	]
func (rh *RepositorieHandler) UpdateManyMetrics(w http.ResponseWriter, r *http.Request) {
	log := rh.Logger.LogrusLog

	if rh.trustIPNet != nil {
		ipStr := r.Header.Get("X-Real-IP")

		ip := net.ParseIP(ipStr)

		allowed := rh.trustIPNet.Contains(ip)

		if !allowed {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	metricsList := []metric.Metric{}
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&metricsList); err != nil {
		log.Errorf("handler func UpdateManyMetrics(): error decode metrics - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}

	err := rh.Repo.UpdateManyMetrics(r.Context(), metricsList)
	if err != nil {
		log.Errorf("handler func UpdateManyMetrics(): error update metrics: %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
