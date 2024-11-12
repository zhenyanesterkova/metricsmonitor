package memstorage

import "github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"

type Memento struct {
	metrics map[string]metric.Metric
}

func (m *Memento) GetSavedState() map[string]metric.Metric {
	return m.metrics
}
