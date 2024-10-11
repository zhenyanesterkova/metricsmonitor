package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric/metricerrors"
)

func (rh *RepositorieHandler) UpdateMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		metricType := chi.URLParam(r, "typeMetric")
		metricName := chi.URLParam(r, "nameMetric")
		metricValue := chi.URLParam(r, "valueMetric")

		err := rh.Repo.UpdateMetric(metricName, metricType, metricValue)
		if err != nil {
			switch err {
			case metricerrors.ErrInvalidName:
				w.WriteHeader(http.StatusNotFound)
				return
			case metricerrors.ErrParseValue, metricerrors.ErrUnknownType, metricerrors.ErrInvalidType:
				w.WriteHeader(http.StatusBadRequest)
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
		}
	}
}
