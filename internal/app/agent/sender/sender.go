package sender

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
)

type Sender struct {
	report                  ReportData
	client                  *http.Client
	hashKey                 *string
	endpoint                string
	requestAttemptIntervals []string
	reportInterval          time.Duration
	rateLimit               int
}

type ReportData struct {
	metricsBuf *metric.MetricBuf
}

func New(
	addr string,
	reportInt time.Duration,
	buff *metric.MetricBuf,
	hashKey *string,
	rateLimit int,
) *Sender {
	return &Sender{
		client:         &http.Client{},
		endpoint:       addr,
		reportInterval: reportInt,
		report: ReportData{
			metricsBuf: buff,
		},
		requestAttemptIntervals: []string{
			"1s",
			"3s",
			"5s",
		},
		hashKey:   hashKey,
		rateLimit: rateLimit,
	}
}

func (s *Sender) SendQueryUpdateMetrics() error {
	mList := s.report.metricsBuf.GetMetricsList()

	if len(mList) == 0 {
		log.Println("    no data for sending ...")
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

	url := s.endpoint

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

	if s.hashKey != nil {
		h := hmac.New(sha256.New, []byte(*s.hashKey))
		h.Write(buff.Bytes())
		sum := hex.EncodeToString(h.Sum(nil))
		req.Header.Set("HashSHA256", sum)
	}

	log.Println("send request ...")
	resp, err := s.client.Do(req)
	defer func(err error) {
		if err == nil {
			errBodyClose := resp.Body.Close()
			if errBodyClose != nil {
				log.Fatalf("sender.go func SendQueryUpdateMetrics(): error close body - %v", errBodyClose)
			}
		}
	}(err)
	if err != nil {
		reqSuccess := false
		for i, interval := range s.requestAttemptIntervals {
			dur, errParse := time.ParseDuration(interval)
			if errParse != nil {
				return fmt.Errorf(`failed send statistic to server: %w;
				the attempt to re-send № %d failed: 
				the interval could not be parsed: %w`,
					err,
					i+1,
					errParse,
				)
			}
			time.Sleep(dur)
			resp, err = s.client.Do(req)
			if err == nil {
				reqSuccess = true
				break
			}
		}
		if !reqSuccess {
			return fmt.Errorf(`failed send statistic to server: %w,
			all attempts to re-send failed`,
				err,
			)
		}
	}
	return nil
}

func (s *Sender) SendReport(ctx context.Context) {
	jobs := make(chan struct{}, 1)

	defer close(jobs)

	for w := 1; w <= s.rateLimit; w++ {
		go s.sendWorker(jobs)
	}
	ticker := time.NewTicker(s.reportInterval)
	for range ticker.C {
		select {
		case <-ctx.Done():
			log.Println("Stop send workers.")
			return
		case jobs <- struct{}{}:
			log.Println("Start send statistic ...")
		}
	}
}

func (s *Sender) sendWorker(
	jobs <-chan struct{},
) {
	for range jobs {
		err := s.SendQueryUpdateMetrics()
		if err != nil {
			log.Printf("error occurred while sending metrics to server %v", err)
			continue
		}
		s.report.metricsBuf.ResetCountersValues()
	}
}
