package sender

import (
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

	url := fmt.Sprintf("http://%s/update/%s/%s/%s", s.Endpoint, upMetric.GetType(), upMetric.GetName(), upMetric.GetValue())

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (s Sender) SendReport() {

	defer s.Report.WGroup.Done()
	ticker := time.NewTicker(s.ReportInterval)
	for range ticker.C {

		for name, _ := range s.Report.MetricsBuf.Metrics {
			s.Report.Mutex.Lock()
			err := s.SendQueryUpdateMetric(name)
			if err != nil {
				log.Printf("an error occurred while sending the report to the server %v", err)
			}
			s.Report.Mutex.Unlock()
		}
	}
}
