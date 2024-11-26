package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
)

const (
	countCheckServer = 3
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

	log.Printf("new request to url=%s, method=%s, data: %s", url, http.MethodPost, upMetric)

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

func (s Sender) SendReport() error {
	defer s.Report.WGroup.Done()

	var ok bool
	for range countCheckServer {
		ok = s.checkServerAvailability()
		if !ok {
			time.Sleep(time.Minute)
			continue
		}
		break
	}

	if !ok {
		return fmt.Errorf("sender.go func SendReport(): %w", errors.New("the server is not available"))
	}

	log.Println("Start report statistic ...")
	ticker := time.NewTicker(s.ReportInterval)
	for range ticker.C {
		s.Report.MetricsBuf.Lock()
		for name := range s.Report.MetricsBuf.Metrics {
			err := s.SendQueryUpdateMetric(name)
			if err != nil {
				log.Printf("an error occurred while sending the report to the server %v", err)
				return fmt.Errorf("sender.go func SendReport(): send metric error - %w", err)
			}
		}
		s.Report.MetricsBuf.Unlock()
		s.Report.MetricsBuf.ResetCountersValues()
	}
	return nil
}

func (s Sender) checkServerAvailability() bool {
	log.Println("Check server availability ...")
	url := fmt.Sprintf("http://%s/", s.Endpoint)
	resp, err := s.Client.Get(url)
	if err != nil {
		log.Println("Connection is not established, waiting 1 minute")
		return false
	}

	err = resp.Body.Close()
	if err != nil {
		return false
	}

	log.Println("Connection is established")
	return true
}
