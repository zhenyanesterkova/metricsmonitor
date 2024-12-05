package statistic

import (
	"sync"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
)

type Statistic struct {
	wGroup       *sync.WaitGroup
	metricsBuf   *metric.MetricBuf
	pollInterval time.Duration
}

func New(
	wg *sync.WaitGroup,
	buff *metric.MetricBuf,
	pInt time.Duration,
) Statistic {
	return Statistic{
		wGroup:       wg,
		metricsBuf:   buff,
		pollInterval: pInt,
	}
}

func (s Statistic) UpdateStatistic() {
	defer s.wGroup.Done()

	ticker := time.NewTicker(s.pollInterval)
	for range ticker.C {
		s.metricsBuf.UpdateMetrics()
	}
}
