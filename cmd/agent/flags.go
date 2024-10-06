package main

import (
	"flag"
	"os"
	"strconv"
)

var (
	endpoint       string
	pollInterval   int
	reportInterval int
)

func parseFlags() {

	flag.StringVar(&endpoint, "a", "localhost:8080", "address and port to send report on server")
	flag.IntVar(&pollInterval, "p", 2, "the frequency of polling metrics from the runtime package")
	flag.IntVar(&reportInterval, "r", 10, "the frequency of sending metrics to the server")

	flag.Parse()

	if envEndpoint := os.Getenv("ADDRESS"); envEndpoint != "" {
		endpoint = envEndpoint
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		pollIntervalTemp, err := strconv.Atoi(envPollInterval)
		if err != nil {
			panic(err)
		}
		pollInterval = pollIntervalTemp
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		envReportIntervalTemp, err := strconv.Atoi(envReportInterval)
		if err != nil {
			panic(err)
		}
		reportInterval = envReportIntervalTemp
	}
}
