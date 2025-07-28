// An agent (HTTP client) for collecting runtime metrics and then sending them to the server over the HTTP protocol.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/grpcsender"
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

	s, err := sender.New(address, cfg.ReportInterval, metrics, cfg.HashKey, cfg.RateLimit, cfg.CryptoKeyPath)
	if err != nil {
		log.Fatalf("an error occurred while create the sender: %v", err)
	}

	sGRPC, err := grpcsender.New(cfg.Address, cfg.ReportInterval, metrics, cfg.HashKey, cfg.RateLimit, cfg.CryptoKeyPath)
	if err != nil {
		log.Fatalf("an error occurred while create the grpc sender: %v", err)
	}
	defer func() {
		err := sGRPC.Close()
		if err != nil {
			log.Fatalf("failed close grpc conn: %v", err)
		}
	}()

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
	go s.SendReport(sendCtx)

	GRPCSendCtx := context.WithoutCancel(ctx)
	go sGRPC.SendReport(GRPCSendCtx)

	select {
	case <-ctx.Done():
		log.Println("Got stop signal")
	case err := <-errCh:
		stop()
		log.Printf("fatal error: %v", err)
	}
}
