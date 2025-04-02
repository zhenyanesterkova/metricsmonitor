package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
	"github.com/zhenyanesterkova/metricsmonitor/web"
)

// GetAllMetrics handles HTTP requests to retrieve all available metrics
// and renders them using an HTML template.
//
// HTTP Response:
//
// - Returns an HTML page listing all metrics
//
// - Sets appropriate HTTP status codes:
//   - 200 OK: Successful retrieval
//   - 500 Internal Server Error: Error
//
// Example usage in HTTP request:
//
//	GET /
func (rh *RepositorieHandler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	log := rh.Logger.LogrusLog

	res, err := rh.Repo.GetAllMetrics()

	if err != nil {
		log.Errorf("handler func GetAllMetrics(): error get metrics - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}

	tmplMetrics, err := template.ParseFS(web.Templates, "template/allMetricsView.html")
	if err != nil {
		log.Errorf("handler func GetAllMetrics(): error parse html template - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	err = tmplMetrics.ExecuteTemplate(w, "metrics", res)
	if err != nil {
		log.Errorf("handler func GetAllMetrics(): error execute html template - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}
}

// GetMetricValue handles HTTP requests to retrieve metric value.
//
// URL Parameters:
//   - nameMetric: The name of the metric to retrieve.
//   - typeMetric: The type of the metric (gauge, counter).
//
// HTTP Response:
//
// - Returns the metric value as a string.
//
// - Sets appropriate HTTP status codes:
//   - 200 OK: Successful retrieval.
//   - 404 Not Found: Unknown metric or invalid type.
//   - 500 Internal Server Error: Other errors.
//
// Example usage in HTTP request:
//
//	GET /value/gauge/Alloc
//
// Example response:
//
//	22.5
func (rh *RepositorieHandler) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	log := rh.Logger.LogrusLog

	name := chi.URLParam(r, "nameMetric")
	metricType := chi.URLParam(r, "typeMetric")

	res, err := rh.Repo.GetMetricValue(name, metricType)

	if err != nil {
		if errors.Is(err, metric.ErrUnknownMetric) || errors.Is(err, metric.ErrInvalidType) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Errorf("handler func GetMetricValue(): error get metric value - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(res.String()))
}

// GetMetricValueJSON handles HTTP requests to retrieve metric value in JSON format.
//
// HTTP Request:
//
// - Accepts a JSON payload with metric details.
//
// - Example request body:
//
//	json
//
//	{
//	  "ID": "metric_name",
//	  "MType": "metric_type"
//	}
//
// HTTP Response:
//
// - Returns the metric value in JSON format.
//
// - Sets appropriate HTTP status codes:
//   - 200 OK: Successful retrieval.
//   - 404 Not Found: Unknown metric or invalid type.
//   - 500 Internal Server Error: Other errors.
//
// Example usage in HTTP request:
//
//	POST /value/
//	Content-Type: application/json
//
//	{
//	  "ID": "Alloc",
//	  "MType": "gauge"
//	}
//
// Example response:
//
//	json
//	{
//	  "ID": "Alloc",
//	  "MType": "gauge",
//	  "value": 22.5
//	}
func (rh *RepositorieHandler) GetMetricValueJSON(w http.ResponseWriter, r *http.Request) {
	log := rh.Logger.LogrusLog

	metrica := metric.New("")
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metrica); err != nil {
		log.Errorf("handler func GetMetricValueJSON(): error decode metric - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}

	res, err := rh.Repo.GetMetricValue(metrica.ID, metrica.MType)

	if err != nil {
		if errors.Is(err, metric.ErrUnknownMetric) || errors.Is(err, metric.ErrInvalidType) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Errorf("handler func GetMetricValueJSON(): error get metric value - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(res); err != nil {
		log.Errorf("handler func GetMetricValueJSON(): error encode metric - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}
}

// Ping checks the availability of the service.
// Returns: status 200 if the service is available,
// status 500 if the storage is not available.
func (rh *RepositorieHandler) Ping(w http.ResponseWriter, r *http.Request) {
	log := rh.Logger.LogrusLog

	err := rh.Repo.Ping()

	if err != nil {
		log.Errorf("failed ping storage: %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
