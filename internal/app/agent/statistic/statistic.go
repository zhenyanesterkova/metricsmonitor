package statistic

import (
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
		s.metricsBuf.UpdateMetrics()
	}
}
