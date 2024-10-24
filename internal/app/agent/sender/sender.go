package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
)

type Sender struct {
	Client         *http.Client
	Endpoint       string
	ReportInterval time.Duration
	Report         ReportData
}

type ReportData struct {
	MetricsBuf *metric.MetricBuf
	WGroup     *sync.WaitGroup
	Mutex      *sync.Mutex
}

func (s Sender) SendQueryUpdateMetric(metricName string) error {

	upMetric := s.Report.MetricsBuf.Metrics[metricName]

	var buff bytes.Buffer
	enc := json.NewEncoder(&buff)
	if err := enc.Encode(upMetric); err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/update/", s.Endpoint)

	req, err := http.NewRequest(http.MethodPost, url, &buff)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}

	responseMetric := metric.Metric{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&responseMetric); err != nil {
		return err
	}

	defer resp.Body.Close()

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
