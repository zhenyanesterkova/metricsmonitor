package handlers

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric/metricerrors"

	"github.com/go-chi/chi/v5"
)

func (rh *RepositorieHandler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {

	res, err := rh.Repo.GetAllMetrics()
	if err != nil {
		http.Error(w, "error get metrics: "+err.Error(), http.StatusInternalServerError)
	}

	templatePath := filepath.Join("../../", "web", "template", "allMetricsView.html")
	tmplIndex, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	err = tmplIndex.ExecuteTemplate(w, "metrics", res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

}
func (rh *RepositorieHandler) GetMetricValue(w http.ResponseWriter, r *http.Request) {

	name := chi.URLParam(r, "nameMetric")
	metricType := chi.URLParam(r, "typeMetric")
	res, err := rh.Repo.GetMetricValue(name, metricType)
	if err != nil {
		if err == metricerrors.ErrUnknownMetric || err == metricerrors.ErrInvalidType {
			w.WriteHeader(http.StatusNotFound)
		}
	}
	_, _ = w.Write([]byte(res))

}
