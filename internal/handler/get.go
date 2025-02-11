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
