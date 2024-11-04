package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/sender"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/statistic"
)

func main() {

	var mutex sync.Mutex
	var wg sync.WaitGroup

	cfg := config.New()
	err := cfg.Build()
	if err != nil {
		log.Fatalf("an error occurred while reading the config %v", err)
	}

	metrics := metric.NewMetricBuf()
	stats := statistic.Statistic{
		PollInterval: cfg.PollInterval,
		Mutex:        &mutex,
		WGroup:       &wg,
		MetricsBuf:   metrics,
	}
	sender := sender.Sender{
		Client:         &http.Client{},
		Endpoint:       cfg.Address,
		ReportInterval: cfg.ReportInterval,
		Report: sender.ReportData{
			MetricsBuf: metrics,
			WGroup:     &wg,
			Mutex:      &mutex,
		},
	}

	wg.Add(2)

	go stats.UpdateStatistic()

	go sender.SendReport()

	wg.Wait()

}
