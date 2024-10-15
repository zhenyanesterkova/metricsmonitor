package handlers

import (
	"net/http"
	"text/template"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric/metricerrors"
	"github.com/zhenyanesterkova/metricsmonitor/web"

	"github.com/go-chi/chi/v5"
)

func (rh *RepositorieHandler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {

	res, err := rh.Repo.GetAllMetrics()
	if err != nil {
		http.Error(w, "error get metrics: "+err.Error(), http.StatusInternalServerError)
	}

	tmplMetrics, err := template.ParseFS(web.Templates, "template/allMetricsView.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	err = tmplMetrics.ExecuteTemplate(w, "metrics", res)
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
