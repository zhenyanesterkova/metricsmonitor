package update

import (
	"net/http"
	"strings"

	"github.com/zhenyanesterkova/metricsmonitor/internal/metric/metricErrors"
)

type Storage interface {
	Update(name string, typeMetric string, val string) error
	String() string
}

func New(s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		params := strings.Split(r.Form.Get("data"), "/")[1:]
		if len(params) != 4 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		metricType := params[1]
		metricName := params[2]
		metricValue := params[3]

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
