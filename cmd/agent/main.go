package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/config"
)

var mutex sync.Mutex
var wg sync.WaitGroup

type metric struct {
	metricType string
	value      string
}

var Metrics = map[string]*metric{
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

func updateMetrics(metrics map[string]*metric, statStruct *runtime.MemStats, mutex *sync.Mutex) error {
	mutex.Lock()
	runtime.ReadMemStats(statStruct)
	statStructFields := reflect.ValueOf(statStruct).Elem()

	for i := 0; i < statStructFields.NumField(); i++ {

		field := statStructFields.Field(i)
		fieldName := statStructFields.Type().Field(i).Name
		fieldType := field.Kind()
		if metric, ok := metrics[fieldName]; ok {
			switch fieldType {
			case reflect.Float64:
				metric.value = strconv.FormatFloat(field.Float(), 'f', -1, 64)
			case reflect.Uint64:
				metric.value = strconv.FormatUint(field.Uint(), 10)
			case reflect.Uint32:
				metric.value = strconv.FormatUint(field.Uint(), 10)
			}
		}

	}

	pollCountValue, err := strconv.ParseInt(metrics["PollCount"].value, 10, 64)
	if err != nil {
		return err
	}
	metrics["PollCount"].value = strconv.FormatInt(pollCountValue+1, 10)

	metrics["RandomValue"].value = strconv.FormatFloat(rand.Float64(), 'f', -1, 64)
	mutex.Unlock()
	return nil
}

func sendQueryUpdateMetric(client *http.Client, mName string, m metric, endpoint string) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%s", endpoint, m.metricType, mName, m.value)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func updateStatistic(interval time.Duration, mutex *sync.Mutex) {
	defer wg.Done()
	stats := &runtime.MemStats{}

	ticker := time.NewTicker(interval)
	for range ticker.C {

		err := updateMetrics(Metrics, stats, mutex)
		if err != nil {
			log.Printf("an error occurred while updating the metrics %v", err)
		}
	}
}

func sendReport(client *http.Client, endpoint string, interval time.Duration, mutex *sync.Mutex) {

	defer wg.Done()
	ticker := time.NewTicker(interval)
	for range ticker.C {

		for name, metric := range Metrics {
			mutex.Lock()
			err := sendQueryUpdateMetric(client, name, *metric, endpoint)
			if err != nil {
				log.Printf("an error occurred while sending the report to the server %v", err)
			}
			mutex.Unlock()
		}
	}
}

func getConfig() (config.Config, error) {
	cfgBuilder := config.GetConfigBuilder()
	cfgDirector := config.NewConfigDirector(cfgBuilder)
	resConfig, err := cfgDirector.BuildConfig()
	if err != nil {
		return resConfig, err
	}

	return resConfig, nil
}

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatalf("an error occurred while reading the config %v", err)
	}

	client := &http.Client{}

	wg.Add(2)
	go updateStatistic(cfg.PollInterval, &mutex)

	go sendReport(client, cfg.Address, cfg.ReportInterval, &mutex)

	wg.Wait()

}
