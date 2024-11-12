package statistic

import (
	"log"
	"sync"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
)

type Statistic struct {
	Mutex        *sync.Mutex
	WGroup       *sync.WaitGroup
	MetricsBuf   *metric.MetricBuf
	PollInterval time.Duration
}

func (s Statistic) UpdateStatistic() {
	defer s.WGroup.Done()

	ticker := time.NewTicker(s.PollInterval)
	for range ticker.C {
		err := s.MetricsBuf.UpdateMetrics(s.Mutex)
		if err != nil {
			log.Printf("error in updating the metrics: %v", err.Error())
		}
	}
}
