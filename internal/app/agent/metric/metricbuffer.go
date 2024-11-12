package metric

import (
	"math/rand/v2"
	"runtime"
	"sync"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type MetricBuf struct {
	Metrics map[string]*Metric
}

func NewMetricBuf() *MetricBuf {
	buffer := &MetricBuf{}
	buffer.Metrics = map[string]*Metric{
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
	}
	return buffer
}

func (buf *MetricBuf) SetGaugeValuesInMetrics() error {
	for _, metric := range buf.Metrics {
		if metric.MType == GaugeType {
			err := metric.setGaugeValue()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (buf *MetricBuf) UpdateMetrics(mutex *sync.Mutex) error {
	statStruct := &runtime.MemStats{}
	runtime.ReadMemStats(statStruct)

	mutex.Lock()

	buf.Metrics["Alloc"].updateGauge(statStruct.Alloc)
	buf.Metrics["BuckHashSys"].updateGauge(statStruct.BuckHashSys)
	buf.Metrics["Frees"].updateGauge(statStruct.Frees)
	buf.Metrics["GCCPUFraction"].updateGauge(statStruct.GCCPUFraction)
	buf.Metrics["GCSys"].updateGauge(statStruct.GCSys)
	buf.Metrics["HeapAlloc"].updateGauge(statStruct.HeapAlloc)
	buf.Metrics["HeapIdle"].updateGauge(statStruct.HeapIdle)
	buf.Metrics["HeapInuse"].updateGauge(statStruct.HeapInuse)
	buf.Metrics["HeapObjects"].updateGauge(statStruct.HeapObjects)
	buf.Metrics["HeapReleased"].updateGauge(statStruct.HeapReleased)
	buf.Metrics["HeapSys"].updateGauge(statStruct.HeapSys)
	buf.Metrics["LastGC"].updateGauge(statStruct.LastGC)
	buf.Metrics["Lookups"].updateGauge(statStruct.Lookups)
	buf.Metrics["MCacheInuse"].updateGauge(statStruct.MCacheInuse)
	buf.Metrics["MCacheSys"].updateGauge(statStruct.MCacheSys)
	buf.Metrics["MSpanInuse"].updateGauge(statStruct.MSpanInuse)
	buf.Metrics["MSpanSys"].updateGauge(statStruct.MSpanSys)
	buf.Metrics["Mallocs"].updateGauge(statStruct.Mallocs)
	buf.Metrics["NextGC"].updateGauge(statStruct.NextGC)
	buf.Metrics["NumForcedGC"].updateGauge(statStruct.NumForcedGC)
	buf.Metrics["NumGC"].updateGauge(statStruct.NumGC)
	buf.Metrics["OtherSys"].updateGauge(statStruct.OtherSys)
	buf.Metrics["PauseTotalNs"].updateGauge(statStruct.PauseTotalNs)
	buf.Metrics["StackInuse"].updateGauge(statStruct.StackInuse)
	buf.Metrics["StackSys"].updateGauge(statStruct.StackSys)
	buf.Metrics["Sys"].updateGauge(statStruct.Sys)
	buf.Metrics["TotalAlloc"].updateGauge(statStruct.TotalAlloc)

	buf.Metrics["PollCount"].updateCounter()

	buf.Metrics["RandomValue"].updateGauge(rand.Float64())

	mutex.Unlock()

	return nil
}
