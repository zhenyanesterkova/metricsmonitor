package main

import (
	"runtime"
	"sync"
	"testing"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"

	"github.com/stretchr/testify/assert"
)

var metricsTest = metric.NewMetricBuf()

func Test_updateMetrics(t *testing.T) {

	type args struct {
		metrics    metric.MetricBuf
		statStruct *runtime.MemStats
		mutex      *sync.Mutex
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test #1 correct update metrics",
			args: args{
				metrics:    *metricsTest,
				statStruct: &runtime.MemStats{},
				mutex:      &sync.Mutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.metrics.UpdateMetrics(tt.args.mutex)

			assert.NoError(t, err)
		})
	}
}
