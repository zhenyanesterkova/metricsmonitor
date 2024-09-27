package memStorage

import (
	"fmt"
	"strings"

	"github.com/zhenyanesterkova/metricsmonitor/internal/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/metric/metricErrors"
)

type Storage struct {
	Metrics map[string]metric.Metric
}

func New() *Storage {
	return &Storage{
		Metrics: make(map[string]metric.Metric),
	}
}

func (s *Storage) Update(name, typeMetric string, val string) error {
	if name == "" {
		return metricErrors.ErrInvalidName
	}

	if curMetric, ok := s.Metrics[name]; ok {
		if curMetric.GetType() != typeMetric {
			return metricErrors.ErrInvalidType
		}
	} else {
		newMetric, err := metric.New(typeMetric)
		if err != nil {
			return err
		}
		s.Metrics[name] = newMetric
	}

	err := s.Metrics[name].SetValue(val)

	return err
}

func (s Storage) String() string {
	var sb strings.Builder
	for k, v := range s.Metrics {
		sb.WriteString(fmt.Sprintf("%s : %s\n", k, v.String()))
	}

	return sb.String()
}
