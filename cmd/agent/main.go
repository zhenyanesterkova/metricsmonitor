package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const pollInterval time.Duration = 2 * time.Second
const reportInterval time.Duration = 10 * time.Second

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

			log.Printf("    %s : %s\n", fieldName, metrics[fieldName].value)
		}

	}

	pollCountValue, err := strconv.ParseInt(metrics["PollCount"].value, 10, 64)
	if err != nil {
		return err
	}
	metrics["PollCount"].value = strconv.FormatInt(pollCountValue+1, 10)
	log.Printf("    PollCount : %v\n", metrics["PollCount"].value)

	metrics["RandomValue"].value = strconv.FormatFloat(rand.Float64(), 'f', -1, 64)
	log.Printf("    RandomValue : %v\n", metrics["RandomValue"].value)
	mutex.Unlock()
	return nil
}

func sendQueryUpdateMetric(client *http.Client, mName string, m metric, endpoint string) error {

	builder := strings.Builder{}
	builder.WriteString(endpoint)
	builder.WriteString(m.metricType)
	builder.WriteString("/")
	builder.WriteString(mName)
	builder.WriteString("/")
	builder.WriteString(m.value)
	url := strings.TrimSpace(builder.String())

	defer builder.Reset()

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Printf("can not create request: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "text/plain")

	log.Printf("send request to %s\n", url)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("can not do request: %v", err)
		return err
	}
	log.Printf("got response from %s\n", url)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Printf("can not do request: %v", err)
		return err
	}
	log.Printf("response body:\n")
	log.Println(string(body))

	return nil
}

func updateStatistic(interval time.Duration, mutex *sync.Mutex) {
	defer wg.Done()
	stats := &runtime.MemStats{}

	ticker := time.NewTicker(interval)
	for range ticker.C {

		err := updateMetrics(Metrics, stats, mutex)
		if err != nil {
			panic(err)
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
				panic(err)
			}
			mutex.Unlock()
		}
	}
}

func main() {
	endpoint := "http://localhost:8080/update/"
	client := &http.Client{}

	wg.Add(2)
	go updateStatistic(pollInterval, &mutex)

	go sendReport(client, endpoint, reportInterval, &mutex)

	wg.Wait()

}
