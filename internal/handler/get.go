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
	res, err := rh.Repo.GetAllMetrics()
	if err != nil {
		http.Error(w, "error get metrics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmplMetrics, err := template.ParseFS(web.Templates, "template/allMetricsView.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	err = tmplMetrics.ExecuteTemplate(w, "metrics", res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rh *RepositorieHandler) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "nameMetric")
	metricType := chi.URLParam(r, "typeMetric")

	res, err := rh.Repo.GetMetricValue(name, metricType)
	if err != nil {
		if errors.Is(err, metric.ErrUnknownMetric) || errors.Is(err, metric.ErrInvalidType) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(res.String()))
}

func (rh *RepositorieHandler) GetMetricValueJSON(w http.ResponseWriter, r *http.Request) {
	metrica := metric.New("")
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metrica); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	res, err := rh.Repo.GetMetricValue(metrica.ID, metrica.MType)
	if err != nil {
		if errors.Is(err, metric.ErrUnknownMetric) || errors.Is(err, metric.ErrInvalidType) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
}
