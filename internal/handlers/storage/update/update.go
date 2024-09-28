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

		pathArr := strings.Split(r.URL.Path, "/")[1:]
		if len(pathArr) != 4 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		metricType := pathArr[1]
		metricName := pathArr[2]
		metricValue := pathArr[3]

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
