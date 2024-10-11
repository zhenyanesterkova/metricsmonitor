package handlers

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric/metricerrors"
)

func (rh *RepositorieHandler) GetAllMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := rh.Repo.GetAllMetrics()
		if err != nil {
			http.Error(w, "error get metrics: "+err.Error(), http.StatusInternalServerError)
		}

		index := filepath.Join("../../", "web", "template", "allMetricsView.html")
		tmplIndex, err := template.ParseFiles(index)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmplIndex.ExecuteTemplate(w, "metrics", res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
func (rh *RepositorieHandler) GetMetricValue() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "nameMetric")
		metricType := chi.URLParam(r, "typeMetric")
		res, err := rh.Repo.GetMetricValue(name, metricType)
		if err != nil {
			if err == metricerrors.ErrUnknownMetric || err == metricerrors.ErrInvalidType {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
		_, _ = w.Write([]byte(res))
	}
}
