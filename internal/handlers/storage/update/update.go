package update

import (
	"net/http"

	"github.com/zhenyanesterkova/metricsmonitor/internal/handlers"
	"github.com/zhenyanesterkova/metricsmonitor/internal/metric/metricerrors"
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
			case metricerrors.ErrInvalidName:
				w.WriteHeader(http.StatusNotFound)
				return
			case metricerrors.ErrInvalidType, metricerrors.ErrParseValue, metricerrors.ErrUnknownType:
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
