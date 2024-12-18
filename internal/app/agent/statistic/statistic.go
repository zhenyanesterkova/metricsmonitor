package statistic

import (
	"context"
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

func (s Statistic) UpdateStatistic(ctx context.Context) {
	ticker := time.NewTicker(s.pollInterval)
	for {
		select {
		case <-ticker.C:
			log.Println("update statistic")
			s.metricsBuf.UpdateMetrics()
		case <-ctx.Done():
			return
		}
	}
}

func (s Statistic) UpdateGopsutilStatistic(ctx context.Context, errCh chan error) {
	defer close(errCh)
	ticker := time.NewTicker(s.pollInterval)
	for {
		select {
		case <-ticker.C:
			log.Println("update gopsutil statistic")
			err := s.metricsBuf.UpdateGopsutilMetrics()
			if err != nil {
				errCh <- fmt.Errorf("failed update gopsutil metrics: %w", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
