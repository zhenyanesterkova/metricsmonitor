package metric

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateGauge(t *testing.T) {
	metrica := Metric{
		ID:    "test",
		MType: "gauge",
	}
	temp := 3.5
	expectMetrica := Metric{
		ID:    "test",
		MType: "gauge",
		Value: &temp,
	}
	t.Run("Test #1", func(t *testing.T) {
		metrica.updateGauge(3.5)
		assert.Equal(t, expectMetrica, metrica)
	})
}

func Test_updateCounter(t *testing.T) {
	metrica := Metric{
		ID:    "test",
		MType: "counter",
	}
	temp := int64(1)
	expectMetrica := Metric{
		ID:    "test",
		MType: "counter",
		Delta: &temp,
	}
	t.Run("Test #1", func(t *testing.T) {
		metrica.updateCounter()
		assert.Equal(t, expectMetrica, metrica)
	})
}

func Test_UpdateMetrics(t *testing.T) {
	var metricsTest = NewMetricBuf()
	type args struct {
		metrics    *MetricBuf
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

func Test_UpdateGopsutilMetrics(t *testing.T) {
	var metricsTest = NewMetricBuf()
	type args struct {
		metrics    *MetricBuf
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
			err := tt.args.metrics.UpdateGopsutilMetrics()
			assert.NoError(t, err)
		})
	}
}
