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

func (s Sender) SendQueryUpdateMetrics() error {
	mList := make([]metric.Metric, 0)
	for _, m := range s.Report.MetricsBuf.Metrics {
		mList = append(mList, *m)
	}

	if len(mList) == 0 {
		log.Printf("    no data for sending ...")
		return nil
	}

	var buff bytes.Buffer

	gzWriter := gzip.NewWriter(&buff)

	enc := json.NewEncoder(gzWriter)
	if err := enc.Encode(mList); err != nil {
		return fmt.Errorf("sender.go func SendQueryUpdateMetrics(): error encode metrics - %w", err)
	}

	err := gzWriter.Close()
	if err != nil {
		return fmt.Errorf("sender.go func SendQueryUpdateMetrics(): error close gzip.Writer - %w", err)
	}

	url := fmt.Sprintf("http://%s/updates/", s.Endpoint)

	log.Printf("new request to url=%s, method=%s", url, http.MethodPost)
	for _, m := range mList {
		log.Printf("    data: %s", &m)
	}

	req, err := http.NewRequest(http.MethodPost, url, &buff)
	if err != nil {
		return fmt.Errorf("sender.go func SendQueryUpdateMetrics(): error create request - %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("sender.go func SendQueryUpdateMetrics(): error do request - %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("sender.go func SendQueryUpdateMetrics(): error close body - %w", err)
	}

	return nil
}

func (s Sender) SendReport() error {
	defer s.Report.WGroup.Done()

	ticker := time.NewTicker(s.ReportInterval)
	for range ticker.C {
		log.Println("Start send statistic ...")
		err := s.SendQueryUpdateMetrics()
		if err != nil {
			log.Printf("an error occurred while sending the report to the server %v", err)
		}

		s.Report.MetricsBuf.ResetCountersValues()
		log.Println("End send statistic ...")
	}
	return nil
}
