package counter

import (
	"strconv"

	"github.com/zhenyanesterkova/metricsmonitor/internal/metric/metricerrors"
)

type MetricCounter struct {
	valueMetric int64
}

func NewMetricCounter() *MetricCounter {
	return &MetricCounter{}
}

func (m *MetricCounter) SetValue(newValue string) error {
	val, err := strconv.ParseInt(newValue, 0, 64)
	if err != nil {
		return metricerrors.ErrParseValue
	}
	m.valueMetric += val
	return nil
}

func (m *MetricCounter) GetType() string {
	return "counter"
}

func (m *MetricCounter) String() string {
	return strconv.FormatInt(m.valueMetric, 10)
}
