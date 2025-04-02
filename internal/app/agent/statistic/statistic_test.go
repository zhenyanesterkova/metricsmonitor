package statistic

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
)

const defaultPollInt = 2 * time.Second

func TestStatistic_UpdateStatistic(t *testing.T) {
	t.Run("Test #1", func(t *testing.T) {
		s := Statistic{
			metricsBuf:   metric.NewMetricBuf(),
			pollInterval: defaultPollInt,
		}
		ctx, stop := context.WithCancel(context.Background())
		go s.UpdateStatistic(ctx)
		time.Sleep(s.pollInterval * 2)
		stop()
	})
}

func TestStatistic_UpdateGopsutilStatistic(t *testing.T) {
	t.Run("Test #1", func(t *testing.T) {
		s := Statistic{
			metricsBuf:   metric.NewMetricBuf(),
			pollInterval: defaultPollInt,
		}
		errCh := make(chan error)
		ctx, stop := context.WithCancel(context.Background())

		go s.UpdateGopsutilStatistic(ctx, errCh)
		time.Sleep(s.pollInterval * 2)
		stop()
		err := <-errCh
		assert.NoError(t, err)
	})
}
