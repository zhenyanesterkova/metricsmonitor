package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/sender"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/statistic"
)

func main() {
	cfg := config.New()
	err := cfg.Build()
	if err != nil {
		log.Fatalf("an error occurred while reading the config %v", err)
	}

	metrics := metric.NewMetricBuf()
	stats := statistic.New(metrics, cfg.PollInterval)
	senderStat := sender.New(cfg.Address, cfg.ReportInterval, metrics, cfg.HashKey, cfg.RateLimit)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	errCh := make(chan error)

	updateCtx := context.WithoutCancel(ctx)
	go stats.UpdateStatistic(updateCtx)

	updateGopsutilCtx := context.WithoutCancel(ctx)
	go stats.UpdateGopsutilStatistic(updateGopsutilCtx, errCh)

	sendCtx := context.WithoutCancel(ctx)
	go senderStat.SendReport(sendCtx)

	select {
	case <-ctx.Done():
		log.Println("Got stop signal")
	case err := <-errCh:
		stop()
		log.Printf("fatal error: %v", err)
	}
}
