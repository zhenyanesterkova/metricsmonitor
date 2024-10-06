package memstorage

import (
	"fmt"
	"strings"

	"github.com/zhenyanesterkova/metricsmonitor/internal/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/metric/metricerrors"
)

type Storage struct {
	metrics map[string]metric.Metric
}

func New() *Storage {
	return &Storage{
		metrics: make(map[string]metric.Metric),
	}
}

func (s *Storage) GetAllMetrics() ([][2]string, error) {
	res := [][2]string{}
	temp := [2]string{}
	for name, metric := range s.metrics {
		temp[0], temp[1] = name, metric.String()
		res = append(res, temp)
	}

	return res, nil
}

func (s *Storage) GetMetricValue(name, typeMetric string) (string, error) {
	if metric, ok := s.metrics[name]; ok {
		if metric.GetType() == typeMetric {
			return metric.String(), nil
		} else {
			return "", metricerrors.ErrInvalidType
		}
	}
	return "", metricerrors.ErrUnknownMetric
}

func (s *Storage) UpdateMetric(name, typeMetric string, val string) error {
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

func (s Storage) String() string {
	var sb strings.Builder
	for k, v := range s.metrics {
		sb.WriteString(fmt.Sprintf("%s : %s\n", k, v.String()))
	}

	return sb.String()
}
