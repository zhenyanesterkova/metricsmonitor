package metric

import (
	"strconv"
)

type MetricCounter struct {
	Delta *int64 `json:"delta,omitempty"`
}

func NewMetricCounter() MetricCounter {
	var delta int64
	return MetricCounter{
		Delta: &delta,
	}
}

func (m *MetricCounter) SetValue(newValue int64) {
	*(m.Delta) += newValue
}

func (m *MetricCounter) String() string {
	return strconv.FormatInt(*(m.Delta), 10)
}
