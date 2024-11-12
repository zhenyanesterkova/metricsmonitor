package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
)

type Sender struct {
	Report         ReportData
	Client         *http.Client
	Endpoint       string
	ReportInterval time.Duration
}

type ReportData struct {
	MetricsBuf *metric.MetricBuf
	WGroup     *sync.WaitGroup
	Mutex      *sync.Mutex
}

func (s Sender) SendQueryUpdateMetric(metricName string) error {
	upMetric := s.Report.MetricsBuf.Metrics[metricName]

	var buff bytes.Buffer

	gzWriter := gzip.NewWriter(&buff)

	enc := json.NewEncoder(gzWriter)
	if err := enc.Encode(upMetric); err != nil {
		return fmt.Errorf("sender.go func SendQueryUpdateMetric(): error encode metric - %w", err)
	}

	err := gzWriter.Close()
	if err != nil {
		return fmt.Errorf("sender.go func SendQueryUpdateMetric(): error close gzip.Writer - %w", err)
	}

	url := fmt.Sprintf("http://%s/update/", s.Endpoint)

	req, err := http.NewRequest(http.MethodPost, url, &buff)
	if err != nil {
		return fmt.Errorf("sender.go func SendQueryUpdateMetric(): error create request - %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("sender.go func SendQueryUpdateMetric(): error do request - %w", err)
	}

	responseMetric := metric.Metric{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&responseMetric); err != nil {
		return fmt.Errorf("sender.go func SendQueryUpdateMetric(): error decode metric - %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("sender.go func SendQueryUpdateMetric(): error close body - %w", err)
	}

	return nil
}

func (s Sender) SendReport() {
	defer s.Report.WGroup.Done()
	ticker := time.NewTicker(s.ReportInterval)
	for range ticker.C {
		s.Report.Mutex.Lock()
		err := s.Report.MetricsBuf.SetGaugeValuesInMetrics()
		if err != nil {
			log.Printf("an error occurred while preparing the data for sending to the server %v", err)
		}
		for name := range s.Report.MetricsBuf.Metrics {
			err := s.SendQueryUpdateMetric(name)
			if err != nil {
				log.Printf("an error occurred while sending the report to the server %v", err)
			}
		}
		s.Report.Mutex.Unlock()
	}
}
