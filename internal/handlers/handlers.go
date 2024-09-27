package handlers

// import (
// 	"encoding/json"
// 	"net/http"
// 	"slices"
// 	"strings"

// 	"github.com/zhenyanesterkova/metricsmonitor/internal/metric"
// )

// func UpdateStorageHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		return
// 	}
// 	if err := r.ParseForm(); err != nil {
// 		w.Write([]byte(err.Error()))
// 		return
// 	}

// 	params := strings.Split(r.Form.Get("data"), "/")[1:]
// 	if len(params) != 4 {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	metricType := params[1]
// 	metricName := params[2]
// 	metricValue := params[3]

// 	if ok := slices.Contains(metric.MetricTypes, metricType); !ok {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	if metricName == "" {
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}

// 	res, _ := json.Marshal(params)

// 	w.Write(res)
// }
