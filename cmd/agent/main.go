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
	var wg sync.WaitGroup

	cfg := config.New()
	err := cfg.Build()
	if err != nil {
		log.Fatalf("an error occurred while reading the config %v", err)
	}

	metrics := metric.NewMetricBuf()
	stats := statistic.Statistic{
		PollInterval: cfg.PollInterval,
		WGroup:       &wg,
		MetricsBuf:   metrics,
	}
	senderStat := sender.Sender{
		Client:         &http.Client{},
		Endpoint:       cfg.Address,
		ReportInterval: cfg.ReportInterval,
		Report: sender.ReportData{
			MetricsBuf: metrics,
			WGroup:     &wg,
		},
		RequestAttemptIntervals: []string{
			"1s",
			"3s",
			"5s",
		},
	}

	go stats.UpdateStatistic()
	wg.Add(1)

	go func() {
		err := senderStat.SendReport()
		if err != nil {
			log.Fatalf("an error occurred while send report on server %v", err)
		}
	}()
	wg.Add(1)

	wg.Wait()
}
