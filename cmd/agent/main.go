// An agent (HTTP client) for collecting runtime metrics and then sending them to the server over the HTTP protocol.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/sender"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/statistic"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {
	cfg := config.New()
	err := cfg.Build()
	if err != nil {
		log.Fatalf("an error occurred while reading the config %v", err)
	}

	metrics := metric.NewMetricBuf()
	stats := statistic.New(metrics, cfg.PollInterval)

	address := fmt.Sprintf("http://%s/updates/", cfg.Address)
	senderStat, err := sender.New(address, cfg.ReportInterval, metrics, cfg.HashKey, cfg.RateLimit, cfg.CryptoKeyPath)
	if err != nil {
		log.Fatalf("an error occurred while create the sender: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	errCh := make(chan error)

	log.Printf("Build version: %s\n", buildVersion)
	log.Printf("Build date: %s\n", buildDate)
	log.Printf("Build commit: %s\n", buildCommit)

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
