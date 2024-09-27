package metric

import (
	"github.com/zhenyanesterkova/metricsmonitor/internal/metric/counter"
	"github.com/zhenyanesterkova/metricsmonitor/internal/metric/gauge"
	"github.com/zhenyanesterkova/metricsmonitor/internal/metric/metricErrors"
)

var MetricTypes = []string{"gauge", "counter"}

type Metric interface {
	SetValue(string) error
	GetType() string
	String() string
}

func New(typeMetric string) (Metric, error) {
	switch typeMetric {
	case "gauge":
		return gauge.NewMetricGauge(), nil
	case "counter":
		return counter.NewMetricCounter(), nil
	default:
		return nil, metricErrors.ErrUnknownType
	}
}
