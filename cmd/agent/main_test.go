package main

import (
	"runtime"
	"testing"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
)

var metricsTest = metric.NewMetricBuf()

func Test_updateMetrics(t *testing.T) {
	type args struct {
		metrics    *metric.MetricBuf
		statStruct *runtime.MemStats
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			name: "test #1 correct update metrics",
			args: args{
				metrics:    metricsTest,
				statStruct: &runtime.MemStats{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.metrics.UpdateMetrics()
		})
	}
}
