package main

import (
	"flag"
	"os"
)

var endpoint string

func parseFlags() {
	flag.StringVar(&endpoint, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if envEndpoint := os.Getenv("ADDRESS"); envEndpoint != "" {
		endpoint = envEndpoint
	}
}
