package statistic

import (
	"errors"
	"sync"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
)

type Statistic struct {
	PollInterval time.Duration
	Mutex        *sync.Mutex
	WGroup       *sync.WaitGroup
	MetricsBuf   *metric.MetricBuf
}

func (s Statistic) UpdateStatistic() error {

	defer s.WGroup.Done()

	ticker := time.NewTicker(s.PollInterval)
	for range ticker.C {

		err := s.MetricsBuf.UpdateMetrics(s.Mutex)
		if err != nil {
			return errors.New("error in updating the metrics: %v" + err.Error())
		}
	}

	return nil
}
