package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric/metricerrors"
)

func getAllMetrics(s Repositorie) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := s.GetAllMetrics()
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
func getMetricValue(s Repositorie) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "nameMetric")
		metricType := chi.URLParam(r, "typeMetric")
		res, err := s.GetMetricValue(name, metricType)
		if err != nil {
			if err == metricerrors.ErrUnknownMetric || err == metricerrors.ErrInvalidType {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}

		io.WriteString(w, res)
	}
}
