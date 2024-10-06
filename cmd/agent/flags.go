package main

import "flag"

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
}
