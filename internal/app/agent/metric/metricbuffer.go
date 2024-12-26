package metric

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type MetricBuf struct {
	Metrics map[string]*Metric
	mutex   *sync.Mutex
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
	}
	buffer.mutex = &sync.Mutex{}
	return buffer
}

func (buf *MetricBuf) GetMetricsList() []Metric {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	mList := make([]Metric, 0)
	for _, m := range buf.Metrics {
		mList = append(mList, *m)
	}
	return mList
}

func (buf *MetricBuf) UpdateMetrics() {
	statStruct := &runtime.MemStats{}
	runtime.ReadMemStats(statStruct)

	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	buf.Metrics["Alloc"].updateGauge(float64(statStruct.Alloc))
	buf.Metrics["BuckHashSys"].updateGauge(float64(statStruct.BuckHashSys))
	buf.Metrics["Frees"].updateGauge(float64(statStruct.Frees))
	buf.Metrics["GCCPUFraction"].updateGauge(float64(statStruct.GCCPUFraction))
	buf.Metrics["GCSys"].updateGauge(float64(statStruct.GCSys))
	buf.Metrics["HeapAlloc"].updateGauge(float64(statStruct.HeapAlloc))
	buf.Metrics["HeapIdle"].updateGauge(float64(statStruct.HeapIdle))
	buf.Metrics["HeapInuse"].updateGauge(float64(statStruct.HeapInuse))
	buf.Metrics["HeapObjects"].updateGauge(float64(statStruct.HeapObjects))
	buf.Metrics["HeapReleased"].updateGauge(float64(statStruct.HeapReleased))
	buf.Metrics["HeapSys"].updateGauge(float64(statStruct.HeapSys))
	buf.Metrics["LastGC"].updateGauge(float64(statStruct.LastGC))
	buf.Metrics["Lookups"].updateGauge(float64(statStruct.Lookups))
	buf.Metrics["MCacheInuse"].updateGauge(float64(statStruct.MCacheInuse))
	buf.Metrics["MCacheSys"].updateGauge(float64(statStruct.MCacheSys))
	buf.Metrics["MSpanInuse"].updateGauge(float64(statStruct.MSpanInuse))
	buf.Metrics["MSpanSys"].updateGauge(float64(statStruct.MSpanSys))
	buf.Metrics["Mallocs"].updateGauge(float64(statStruct.Mallocs))
	buf.Metrics["NextGC"].updateGauge(float64(statStruct.NextGC))
	buf.Metrics["NumForcedGC"].updateGauge(float64(statStruct.NumForcedGC))
	buf.Metrics["NumGC"].updateGauge(float64(statStruct.NumGC))
	buf.Metrics["OtherSys"].updateGauge(float64(statStruct.OtherSys))
	buf.Metrics["PauseTotalNs"].updateGauge(float64(statStruct.PauseTotalNs))
	buf.Metrics["StackInuse"].updateGauge(float64(statStruct.StackInuse))
	buf.Metrics["StackSys"].updateGauge(float64(statStruct.StackSys))
	buf.Metrics["Sys"].updateGauge(float64(statStruct.Sys))
	buf.Metrics["TotalAlloc"].updateGauge(float64(statStruct.TotalAlloc))

	buf.Metrics["PollCount"].updateCounter()

	buf.Metrics["RandomValue"].updateGauge(rand.Float64())
}

func (buf *MetricBuf) UpdateGopsutilMetrics() error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed get memory metric values: %w", err)
	}

	cpuCount, err := cpu.Counts(true)
	if err != nil {
		return fmt.Errorf("failed get cpu metric value: %w", err)
	}

	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	buf.Metrics["TotalMemory"].updateGauge(float64(v.Total))
	buf.Metrics["FreeMemory"].updateGauge(float64(v.Free))
	buf.Metrics["CPUutilization1"].updateGauge(float64(cpuCount))

	return nil
}

func (buf *MetricBuf) ResetCountersValues() {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()
	for _, metrica := range buf.Metrics {
		if metrica.MType == CounterType {
			*metrica.Delta = 0
		}
	}
}
