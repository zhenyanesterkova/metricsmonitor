package metric

import (
	"math/rand/v2"
	"runtime"
	"strconv"
	"sync"
)

type Metric struct {
	name       string
	metricType string
	value      string
}

func (m Metric) GetValue() string {
	return m.value
}
func (m Metric) GetType() string {
	return m.metricType
}
func (m Metric) GetName() string {
	return m.name
}

func (m *Metric) update(val string) {
	m.value = val
}

type MetricBuf struct {
	Metrics map[string]*Metric
}

func NewMetricBuf() *MetricBuf {
	buffer := &MetricBuf{}
	buffer.Metrics = map[string]*Metric{
		"Alloc": &Metric{
			name:       "Alloc",
			metricType: "gauge",
		},
		"BuckHashSys": &Metric{
			name:       "BuckHashSys",
			metricType: "gauge",
		},
		"Frees": &Metric{
			name:       "Frees",
			metricType: "gauge",
		},
		"GCCPUFraction": &Metric{
			name:       "GCCPUFraction",
			metricType: "gauge",
		},
		"GCSys": &Metric{
			name:       "GCSys",
			metricType: "gauge",
		},
		"HeapAlloc": &Metric{
			name:       "HeapAlloc",
			metricType: "gauge",
		},
		"HeapIdle": &Metric{
			name:       "HeapIdle",
			metricType: "gauge",
		},
		"HeapInuse": &Metric{
			name:       "HeapInuse",
			metricType: "gauge",
		},
		"HeapObjects": &Metric{
			name:       "HeapObjects",
			metricType: "gauge",
		},
		"HeapReleased": &Metric{
			name:       "HeapReleased",
			metricType: "gauge",
		},
		"HeapSys": &Metric{
			name:       "HeapSys",
			metricType: "gauge",
		},
		"LastGC": &Metric{
			name:       "LastGC",
			metricType: "gauge",
		},
		"Lookups": &Metric{
			name:       "Lookups",
			metricType: "gauge",
		},
		"MCacheInuse": &Metric{
			name:       "MCacheInuse",
			metricType: "gauge",
		},
		"MCacheSys": &Metric{
			name:       "MCacheSys",
			metricType: "gauge",
		},
		"MSpanInuse": &Metric{
			name:       "MSpanInuse",
			metricType: "gauge",
		},
		"MSpanSys": &Metric{
			name:       "MSpanSys",
			metricType: "gauge",
		},
		"Mallocs": &Metric{
			name:       "Mallocs",
			metricType: "gauge",
		},
		"NextGC": &Metric{
			name:       "NextGC",
			metricType: "gauge",
		},
		"NumForcedGC": &Metric{
			name:       "NumForcedGC",
			metricType: "gauge",
		},
		"NumGC": &Metric{
			name:       "NumGC",
			metricType: "gauge",
		},
		"OtherSys": &Metric{
			name:       "OtherSys",
			metricType: "gauge",
		},
		"PauseTotalNs": &Metric{
			name:       "PauseTotalNs",
			metricType: "gauge",
		},
		"StackInuse": &Metric{
			name:       "StackInuse",
			metricType: "gauge",
		},
		"StackSys": &Metric{
			name:       "StackSys",
			metricType: "gauge",
		},
		"Sys": &Metric{
			name:       "Sys",
			metricType: "gauge",
		},
		"TotalAlloc": &Metric{
			name:       "TotalAlloc",
			metricType: "gauge",
		},
		"PollCount": &Metric{
			name:       "PollCount",
			metricType: "counter",
			value:      "0",
		},
		"RandomValue": &Metric{
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

	buf.Metrics["Alloc"].update(strconv.FormatUint(statStruct.Alloc, 10))
	buf.Metrics["BuckHashSys"].update(strconv.FormatUint(statStruct.BuckHashSys, 10))
	buf.Metrics["Frees"].update(strconv.FormatUint(statStruct.Frees, 10))
	buf.Metrics["GCCPUFraction"].update(strconv.FormatFloat(statStruct.GCCPUFraction, 'f', -1, 64))
	buf.Metrics["GCSys"].update(strconv.FormatUint(statStruct.GCSys, 10))
	buf.Metrics["HeapAlloc"].update(strconv.FormatUint(statStruct.HeapAlloc, 10))
	buf.Metrics["HeapIdle"].update(strconv.FormatUint(statStruct.HeapIdle, 10))
	buf.Metrics["HeapInuse"].update(strconv.FormatUint(statStruct.HeapInuse, 10))
	buf.Metrics["HeapObjects"].update(strconv.FormatUint(statStruct.HeapObjects, 10))
	buf.Metrics["HeapReleased"].update(strconv.FormatUint(statStruct.HeapReleased, 10))
	buf.Metrics["HeapSys"].update(strconv.FormatUint(statStruct.HeapSys, 10))
	buf.Metrics["LastGC"].update(strconv.FormatUint(statStruct.LastGC, 10))
	buf.Metrics["Lookups"].update(strconv.FormatUint(statStruct.Lookups, 10))
	buf.Metrics["MCacheInuse"].update(strconv.FormatUint(statStruct.MCacheInuse, 10))
	buf.Metrics["MCacheSys"].update(strconv.FormatUint(statStruct.MCacheSys, 10))
	buf.Metrics["MSpanInuse"].update(strconv.FormatUint(statStruct.MSpanInuse, 10))
	buf.Metrics["MSpanSys"].update(strconv.FormatUint(statStruct.MSpanSys, 10))
	buf.Metrics["Mallocs"].update(strconv.FormatUint(statStruct.Mallocs, 10))
	buf.Metrics["NextGC"].update(strconv.FormatUint(statStruct.NextGC, 10))
	buf.Metrics["NumForcedGC"].update(strconv.FormatUint(uint64(statStruct.NumForcedGC), 10))
	buf.Metrics["NumGC"].update(strconv.FormatUint(uint64(statStruct.NumGC), 10))
	buf.Metrics["OtherSys"].update(strconv.FormatUint(statStruct.OtherSys, 10))
	buf.Metrics["PauseTotalNs"].update(strconv.FormatUint(statStruct.PauseTotalNs, 10))
	buf.Metrics["StackInuse"].update(strconv.FormatUint(statStruct.StackInuse, 10))
	buf.Metrics["StackSys"].update(strconv.FormatUint(statStruct.StackSys, 10))
	buf.Metrics["Sys"].update(strconv.FormatUint(statStruct.Sys, 10))
	buf.Metrics["TotalAlloc"].update(strconv.FormatUint(statStruct.TotalAlloc, 10))

	pollCountValue, err := strconv.ParseInt(buf.Metrics["PollCount"].value, 10, 64)
	if err != nil {
		return err
	}
	buf.Metrics["PollCount"].update(strconv.FormatInt(pollCountValue+1, 10))

	buf.Metrics["RandomValue"].update(strconv.FormatFloat(rand.Float64(), 'f', -1, 64))

	mutex.Unlock()

	return nil
}
