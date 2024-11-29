package metric

import (
	"strconv"
)

type MetricGauge struct {
	Value *float64 `json:"value,omitempty"`
}

func NewMetricGauge() MetricGauge {
	var val float64
	return MetricGauge{
		Value: &val,
	}
}

func (m *MetricGauge) SetValue(newValue float64) {
	*(m.Value) = newValue
}

func (m *MetricGauge) String() string {
	return strconv.FormatFloat(*(m.Value), 'g', -1, 64)
}
