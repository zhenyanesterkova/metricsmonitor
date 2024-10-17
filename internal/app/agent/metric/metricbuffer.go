package metric

import (
	"math/rand/v2"
	"runtime"
	"sync"
)

type MetricBuf struct {
	Metrics map[string]*metric
}

func NewMetricBuf() *MetricBuf {
	buffer := &MetricBuf{}
	buffer.Metrics = map[string]*metric{
		"Alloc": &metric{
			name:       "Alloc",
			metricType: "gauge",
		},
		"BuckHashSys": &metric{
			name:       "BuckHashSys",
			metricType: "gauge",
		},
		"Frees": &metric{
			name:       "Frees",
			metricType: "gauge",
		},
		"GCCPUFraction": &metric{
			name:       "GCCPUFraction",
			metricType: "gauge",
		},
		"GCSys": &metric{
			name:       "GCSys",
			metricType: "gauge",
		},
		"HeapAlloc": &metric{
			name:       "HeapAlloc",
			metricType: "gauge",
		},
		"HeapIdle": &metric{
			name:       "HeapIdle",
			metricType: "gauge",
		},
		"HeapInuse": &metric{
			name:       "HeapInuse",
			metricType: "gauge",
		},
		"HeapObjects": &metric{
			name:       "HeapObjects",
			metricType: "gauge",
		},
		"HeapReleased": &metric{
			name:       "HeapReleased",
			metricType: "gauge",
		},
		"HeapSys": &metric{
			name:       "HeapSys",
			metricType: "gauge",
		},
		"LastGC": &metric{
			name:       "LastGC",
			metricType: "gauge",
		},
		"Lookups": &metric{
			name:       "Lookups",
			metricType: "gauge",
		},
		"MCacheInuse": &metric{
			name:       "MCacheInuse",
			metricType: "gauge",
		},
		"MCacheSys": &metric{
			name:       "MCacheSys",
			metricType: "gauge",
		},
		"MSpanInuse": &metric{
			name:       "MSpanInuse",
			metricType: "gauge",
		},
		"MSpanSys": &metric{
			name:       "MSpanSys",
			metricType: "gauge",
		},
		"Mallocs": &metric{
			name:       "Mallocs",
			metricType: "gauge",
		},
		"NextGC": &metric{
			name:       "NextGC",
			metricType: "gauge",
		},
		"NumForcedGC": &metric{
			name:       "NumForcedGC",
			metricType: "gauge",
		},
		"NumGC": &metric{
			name:       "NumGC",
			metricType: "gauge",
		},
		"OtherSys": &metric{
			name:       "OtherSys",
			metricType: "gauge",
		},
		"PauseTotalNs": &metric{
			name:       "PauseTotalNs",
			metricType: "gauge",
		},
		"StackInuse": &metric{
			name:       "StackInuse",
			metricType: "gauge",
		},
		"StackSys": &metric{
			name:       "StackSys",
			metricType: "gauge",
		},
		"Sys": &metric{
			name:       "Sys",
			metricType: "gauge",
		},
		"TotalAlloc": &metric{
			name:       "TotalAlloc",
			metricType: "gauge",
		},
		"PollCount": &metric{
			name:       "PollCount",
			metricType: "counter",
			value:      "0",
		},
		"RandomValue": &metric{
			name:       "RandomValue",
			metricType: "gauge",
		},
	}
	return buffer
}

func (buf *MetricBuf) UpdateMetrics(mutex *sync.Mutex) error {

	statStruct := &runtime.MemStats{}
	runtime.ReadMemStats(statStruct)

	mutex.Lock()

	buf.Metrics["Alloc"].update(statStruct.Alloc)
	buf.Metrics["BuckHashSys"].update(statStruct.BuckHashSys)
	buf.Metrics["Frees"].update(statStruct.Frees)
	buf.Metrics["GCCPUFraction"].update(statStruct.GCCPUFraction)
	buf.Metrics["GCSys"].update(statStruct.GCSys)
	buf.Metrics["HeapAlloc"].update(statStruct.HeapAlloc)
	buf.Metrics["HeapIdle"].update(statStruct.HeapIdle)
	buf.Metrics["HeapInuse"].update(statStruct.HeapInuse)
	buf.Metrics["HeapObjects"].update(statStruct.HeapObjects)
	buf.Metrics["HeapReleased"].update(statStruct.HeapReleased)
	buf.Metrics["HeapSys"].update(statStruct.HeapSys)
	buf.Metrics["LastGC"].update(statStruct.LastGC)
	buf.Metrics["Lookups"].update(statStruct.Lookups)
	buf.Metrics["MCacheInuse"].update(statStruct.MCacheInuse)
	buf.Metrics["MCacheSys"].update(statStruct.MCacheSys)
	buf.Metrics["MSpanInuse"].update(statStruct.MSpanInuse)
	buf.Metrics["MSpanSys"].update(statStruct.MSpanSys)
	buf.Metrics["Mallocs"].update(statStruct.Mallocs)
	buf.Metrics["NextGC"].update(statStruct.NextGC)
	buf.Metrics["NumForcedGC"].update(statStruct.NumForcedGC)
	buf.Metrics["NumGC"].update(statStruct.NumGC)
	buf.Metrics["OtherSys"].update(statStruct.OtherSys)
	buf.Metrics["PauseTotalNs"].update(statStruct.PauseTotalNs)
	buf.Metrics["StackInuse"].update(statStruct.StackInuse)
	buf.Metrics["StackSys"].update(statStruct.StackSys)
	buf.Metrics["Sys"].update(statStruct.Sys)
	buf.Metrics["TotalAlloc"].update(statStruct.TotalAlloc)

	err := buf.Metrics["PollCount"].updateCounter()
	if err != nil {
		return err
	}

	buf.Metrics["RandomValue"].update(rand.Float64())

	mutex.Unlock()

	return nil
}
