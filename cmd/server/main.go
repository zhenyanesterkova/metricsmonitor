// A server for collecting runtime metrics that collects reports from agents over the HTTP protocol.
package main

import (
	"fmt"
	"log"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/server"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {
	if err := run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func run() error {
	buildInfo := server.BuildInfo{
		Version: buildVersion,
		Date:    buildDate,
		Commit:  buildCommit,
	}

	srv := server.New(server.HTTPServer, buildInfo)

	if err := srv.Initialize(); err != nil {
		return fmt.Errorf("%w", err)
	}
	defer func() {
		if err := srv.Close(); err != nil {
			log.Printf("error closing server: %v", err)
		}
	}()

	if err := srv.Run(); err != nil {
		return fmt.Errorf("failed to run grpc server: %w", err)
	}

	return nil
}
