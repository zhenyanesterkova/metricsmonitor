package memstorage

import (
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
)

type MemStorage struct {
	metrics map[string]metric.Metric
}

func New() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]metric.Metric),
	}
}

func (s *MemStorage) GetAllMetrics() ([][2]string, error) {
	res := [][2]string{}
	temp := [2]string{}
	for name, metric := range s.metrics {
		temp[0], temp[1] = name, metric.String()
		res = append(res, temp)
	}

	return res, nil
}

func (s *MemStorage) GetMetricValue(name, typeMetric string) (metric.Metric, error) {
	if metrica, ok := s.metrics[name]; ok {
		if metrica.GetType() == typeMetric {
			return metrica, nil
		} else {
			return metric.Metric{}, metric.ErrInvalidType
		}
	}
	return metric.Metric{}, metric.ErrUnknownMetric
}

func (s *MemStorage) UpdateMetric(newMetric metric.Metric) (metric.Metric, error) {
	if newMetric.ID == "" {
		return metric.Metric{}, metric.ErrInvalidName
	}
	if newMetric.MType == "" {
		return metric.Metric{}, metric.ErrInvalidType
	}
	if newMetric.Value == nil && newMetric.Delta == nil {
		return metric.Metric{}, metric.ErrParseValue
	}

	curMetric, ok := s.metrics[newMetric.ID]
	if !ok {
		curMetric = metric.New(newMetric.MType)
		curMetric.ID = newMetric.ID
	}
	if curMetric.GetType() != newMetric.MType {
		return metric.Metric{}, metric.ErrInvalidType
	}

	switch curMetric.MType {
	case metric.TypeGauge:
		curMetric.MetricGauge.SetValue(*newMetric.Value)
	case metric.TypeCounter:
		curMetric.MetricCounter.SetValue(*newMetric.Delta)
	}

	s.metrics[newMetric.ID] = curMetric

	return curMetric, nil
}

func (s *MemStorage) CreateMemento() *Memento {
	return &Memento{Metrics: s.metrics}
}

func (s *MemStorage) RestoreMemento(m *Memento) {
	s.metrics = m.GetSavedState()
}
