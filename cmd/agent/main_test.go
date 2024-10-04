package main

import (
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var metricsTest = map[string]*metric{
	"Alloc": &metric{
		metricType: "gauge",
	},
	"BuckHashSys": &metric{
		metricType: "gauge",
	},
	"Frees": &metric{
		metricType: "gauge",
	},
	"GCCPUFraction": &metric{
		metricType: "gauge",
	},
	"GCSys": &metric{
		metricType: "gauge",
	},
	"HeapAlloc": &metric{
		metricType: "gauge",
	},
	"HeapIdle": &metric{
		metricType: "gauge",
	},
	"HeapInuse": &metric{
		metricType: "gauge",
	},
	"HeapObjects": &metric{
		metricType: "gauge",
	},
	"HeapReleased": &metric{
		metricType: "gauge",
	},
	"HeapSys": &metric{
		metricType: "gauge",
	},
	"LastGC": &metric{
		metricType: "gauge",
	},
	"Lookups": &metric{
		metricType: "gauge",
	},
	"MCacheInuse": &metric{
		metricType: "gauge",
	},
	"MCacheSys": &metric{
		metricType: "gauge",
	},
	"MSpanInuse": &metric{
		metricType: "gauge",
	},
	"MSpanSys": &metric{
		metricType: "gauge",
	},
	"Mallocs": &metric{
		metricType: "gauge",
	},
	"NextGC": &metric{
		metricType: "gauge",
	},
	"NumForcedGC": &metric{
		metricType: "gauge",
	},
	"NumGC": &metric{
		metricType: "gauge",
	},
	"OtherSys": &metric{
		metricType: "gauge",
	},
	"PauseTotalNs": &metric{
		metricType: "gauge",
	},
	"StackInuse": &metric{
		metricType: "gauge",
	},
	"StackSys": &metric{
		metricType: "gauge",
	},
	"Sys": &metric{
		metricType: "gauge",
	},
	"TotalAlloc": &metric{
		metricType: "gauge",
	},
	"PollCount": &metric{
		metricType: "counter",
		value:      "0",
	},
	"RandomValue": &metric{
		metricType: "gauge",
	},
}

func Test_updateMetrics(t *testing.T) {

	type args struct {
		metrics    map[string]*metric
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
				metrics:    metricsTest,
				statStruct: &runtime.MemStats{},
				mutex:      &sync.Mutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := updateMetrics(tt.args.metrics, tt.args.statStruct, tt.args.mutex)

			assert.NoError(t, err)
		})
	}
}
