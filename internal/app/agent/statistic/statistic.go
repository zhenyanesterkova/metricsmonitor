package statistic

import (
	"fmt"
	"log"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
)

type Statistic struct {
	metricsBuf   *metric.MetricBuf
	pollInterval time.Duration
}

func New(
	buff *metric.MetricBuf,
	pInt time.Duration,
) Statistic {
	return Statistic{
		metricsBuf:   buff,
		pollInterval: pInt,
	}
}

func (s Statistic) UpdateStatistic() {
	ticker := time.NewTicker(s.pollInterval)
	for range ticker.C {
		log.Println("update statistic")
		s.metricsBuf.UpdateMetrics()
	}
}

func (s Statistic) UpdateGopsutilStatistic() error {
	ticker := time.NewTicker(s.pollInterval)
	for range ticker.C {
		log.Println("update gopsutil statistic")
		err := s.metricsBuf.UpdateGopsutilMetrics()
		if err != nil {
			return fmt.Errorf("failed update gopsutil metrics: %w", err)
		}
	}
	return nil
}
