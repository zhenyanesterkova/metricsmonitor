package gauge

import (
	"strconv"

	"github.com/zhenyanesterkova/metricsmonitor/internal/metric/metricErrors"
)

type MetricGauge struct {
	valueMetric float64
}

func NewMetricGauge() *MetricGauge {
	return &MetricGauge{}
}

func (m *MetricGauge) SetValue(newValue string) error {
	val, err := strconv.ParseFloat(newValue, 64)
	if err != nil {
		return metricErrors.ErrParseValue
	}
	m.valueMetric = val
	return nil
}

func (m *MetricGauge) GetType() string {
	return "gauge"
}

func (m *MetricGauge) String() string {
	return strconv.FormatFloat(m.valueMetric, 'g', -1, 64)
}
