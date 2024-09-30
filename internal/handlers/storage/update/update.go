package update

import (
	"net/http"

	"github.com/zhenyanesterkova/metricsmonitor/internal/handlers"
	"github.com/zhenyanesterkova/metricsmonitor/internal/metric/metricErrors"
)

func New(s handlers.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		metricType := r.PathValue("typeMetric")
		metricName := r.PathValue("nameMetric")
		metricValue := r.PathValue("valueMetric")

		err := s.Update(metricName, metricType, metricValue)
		if err != nil {
			switch err {
			case metricErrors.ErrInvalidName:
				w.WriteHeader(http.StatusNotFound)
				return
			case metricErrors.ErrInvalidType, metricErrors.ErrParseValue, metricErrors.ErrUnknownType:
				w.WriteHeader(http.StatusBadRequest)
				return
			default:
				w.Write([]byte(err.Error()))
				return
			}
		}

		w.Write([]byte(s.String()))
	}
}
