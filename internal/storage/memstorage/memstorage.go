package memstorage

import (
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric/metricerrors"
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

func (s *MemStorage) GetMetricValue(name, typeMetric string) (string, error) {
	if metric, ok := s.metrics[name]; ok {
		if metric.GetType() == typeMetric {
			return metric.String(), nil
		} else {
			return "", metricerrors.ErrInvalidType
		}
	}
	return "", metricerrors.ErrUnknownMetric
}

func (s *MemStorage) UpdateMetric(name, typeMetric string, val string) error {
	if name == "" {
		return metricerrors.ErrInvalidName
	}
	if typeMetric == "" {
		return metricerrors.ErrInvalidType
	}
	if val == "" {
		return metricerrors.ErrParseValue
	}

	if curMetric, ok := s.metrics[name]; ok {
		if curMetric.GetType() != typeMetric {
			return metricerrors.ErrInvalidType
		}
		err := s.metrics[name].SetValue(val)
		if err != nil {
			return err
		}
	} else {
		newMetric, err := metric.New(typeMetric)
		if err != nil {
			return err
		}
		err = newMetric.SetValue(val)
		if err != nil {
			return err
		}
		s.metrics[name] = newMetric
	}

	return nil
}
