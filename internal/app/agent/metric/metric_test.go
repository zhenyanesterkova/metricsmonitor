package metric

import (
	"sync"
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

func Test_NewMetricBuf(t *testing.T) {
	expBuff := MetricBuf{
		Metrics: map[string]*Metric{
			"Alloc": &Metric{
				ID:    "Alloc",
				MType: "gauge",
			},
			"BuckHashSys": &Metric{
				ID:    "BuckHashSys",
				MType: "gauge",
			},
			"Frees": &Metric{
				ID:    "Frees",
				MType: "gauge",
			},
			"GCCPUFraction": &Metric{
				ID:    "GCCPUFraction",
				MType: "gauge",
			},
			"GCSys": &Metric{
				ID:    "GCSys",
				MType: "gauge",
			},
			"HeapAlloc": &Metric{
				ID:    "HeapAlloc",
				MType: "gauge",
			},
			"HeapIdle": &Metric{
				ID:    "HeapIdle",
				MType: "gauge",
			},
			"HeapInuse": &Metric{
				ID:    "HeapInuse",
				MType: "gauge",
			},
			"HeapObjects": &Metric{
				ID:    "HeapObjects",
				MType: "gauge",
			},
			"HeapReleased": &Metric{
				ID:    "HeapReleased",
				MType: "gauge",
			},
			"HeapSys": &Metric{
				ID:    "HeapSys",
				MType: "gauge",
			},
			"LastGC": &Metric{
				ID:    "LastGC",
				MType: "gauge",
			},
			"Lookups": &Metric{
				ID:    "Lookups",
				MType: "gauge",
			},
			"MCacheInuse": &Metric{
				ID:    "MCacheInuse",
				MType: "gauge",
			},
			"MCacheSys": &Metric{
				ID:    "MCacheSys",
				MType: "gauge",
			},
			"MSpanInuse": &Metric{
				ID:    "MSpanInuse",
				MType: "gauge",
			},
			"MSpanSys": &Metric{
				ID:    "MSpanSys",
				MType: "gauge",
			},
			"Mallocs": &Metric{
				ID:    "Mallocs",
				MType: "gauge",
			},
			"NextGC": &Metric{
				ID:    "NextGC",
				MType: "gauge",
			},
			"NumForcedGC": &Metric{
				ID:    "NumForcedGC",
				MType: "gauge",
			},
			"NumGC": &Metric{
				ID:    "NumGC",
				MType: "gauge",
			},
			"OtherSys": &Metric{
				ID:    "OtherSys",
				MType: "gauge",
			},
			"PauseTotalNs": &Metric{
				ID:    "PauseTotalNs",
				MType: "gauge",
			},
			"StackInuse": &Metric{
				ID:    "StackInuse",
				MType: "gauge",
			},
			"StackSys": &Metric{
				ID:    "StackSys",
				MType: "gauge",
			},
			"Sys": &Metric{
				ID:    "Sys",
				MType: "gauge",
			},
			"TotalAlloc": &Metric{
				ID:    "TotalAlloc",
				MType: "gauge",
			},
			"PollCount": &Metric{
				ID:    "PollCount",
				MType: "counter",
			},
			"RandomValue": &Metric{
				ID:    "RandomValue",
				MType: "gauge",
			},
			"TotalMemory": &Metric{
				ID:    "TotalMemory",
				MType: "gauge",
			},
			"FreeMemory": &Metric{
				ID:    "FreeMemory",
				MType: "gauge",
			},
			"CPUutilization1": &Metric{
				ID:    "CPUutilization1",
				MType: "gauge",
			},
		},
		mutex: &sync.Mutex{},
	}
	t.Run("Test #1", func(t *testing.T) {
		buff := NewMetricBuf()
		assert.Equal(t, expBuff, *buff)
	})
}
