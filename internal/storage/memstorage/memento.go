package memstorage

import "github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"

type Memento struct {
	Metrics map[string]metric.Metric `json:"metrics"`
}

func (m *Memento) GetSavedState() map[string]metric.Metric {
	return m.Metrics
}
