package main

import (
	"fmt"
	"log"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/sender"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/statistic"
)

func main() {
	g := new(errgroup.Group)
	var wg sync.WaitGroup

	cfg := config.New()
	err := cfg.Build()
	if err != nil {
		log.Fatalf("an error occurred while reading the config %v", err)
	}

	metrics := metric.NewMetricBuf()
	stats := statistic.New(metrics, cfg.PollInterval)
	senderStat := sender.New(cfg.Address, cfg.ReportInterval, metrics, cfg.HashKey, cfg.RateLimit)

	go stats.UpdateStatistic()
	wg.Add(1)

	g.Go(func() error {
		err := stats.UpdateGopsutilStatistic()
		if err != nil {
			log.Fatalf("an error occurred while update gopsutil statistic %v", err)
			return fmt.Errorf("an error occurred while update gopsutil statistic %w", err)
		}
		return nil
	})

	go senderStat.SendReport()
	wg.Add(1)

	if err := g.Wait(); err != nil {
		log.Fatalf("fatal error: %v", err)
	}
	wg.Wait()
}
