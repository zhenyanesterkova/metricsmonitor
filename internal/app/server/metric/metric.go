package metric

import (
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric/counter"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric/gauge"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric/metricerrors"
)

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
		return nil, metricerrors.ErrUnknownType
	}
}
