package metric

import (
	"strconv"
)

type Metric struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

func (m *Metric) StringValue() string {
	if m.Delta == nil && m.Value == nil {
		return ""
	}
	if m.Delta == nil {
		return strconv.FormatFloat(*(m.Value), 'f', 2, 64)
	}
	return strconv.FormatInt(*(m.Delta), 10)
}

func (m *Metric) updateGauge(val float64) {
	temp := val
	m.Value = &temp
}

func (m *Metric) updateCounter() {
	if m.Delta == nil {
		temp := int64(0)
		m.Delta = &temp
	}
	*(m.Delta)++
}
