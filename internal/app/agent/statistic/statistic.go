package statistic

import (
	"sync"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
)

type Statistic struct {
	WGroup       *sync.WaitGroup
	MetricsBuf   *metric.MetricBuf
	PollInterval time.Duration
}

func (s Statistic) UpdateStatistic() {
	defer s.WGroup.Done()

	ticker := time.NewTicker(s.PollInterval)
	for range ticker.C {
		s.MetricsBuf.UpdateMetrics()
	}
}
