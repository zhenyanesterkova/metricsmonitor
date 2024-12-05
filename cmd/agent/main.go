package main

import (
	"log"
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
	stats := statistic.New(&wg, metrics, cfg.PollInterval)
	senderStat := sender.New(cfg.Address, cfg.ReportInterval, &wg, metrics)

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
