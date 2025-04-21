package sender

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
)

const (
	aesKeySize = 32
)

type Sender struct {
	report                  ReportData
	client                  *http.Client
	hashKey                 *string
	publicKey               *rsa.PublicKey
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
	pathToPublicKey string,
) (*Sender, error) {
	publicKeyPEM, err := os.ReadFile(pathToPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed read public key from file: %w", err)
	}

	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed parses a public key in PKIX, ASN.1 DER form: %w", err)
	}
	publicKeyRsa, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed converting type to *rsa.PublicKey: %w", err)
	}
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
		publicKey: publicKeyRsa,
	}, nil
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

	aesKey := make([]byte, aesKeySize)
	_, err = rand.Read(aesKey)
	if err != nil {
		return fmt.Errorf("failed to generate AES key: %w", err)
	}

	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	if err != nil {
		return fmt.Errorf("failed to generate IV: %w", err)
	}

	block, _ := aes.NewCipher(aesKey)
	stream := cipher.NewCTR(block, iv)

	ciphertext := make([]byte, len(buff.Bytes()))
	stream.XORKeyStream(ciphertext, buff.Bytes())

	var keyToEncrypt []byte
	keyToEncrypt = append(keyToEncrypt, aesKey...)
	keyToEncrypt = append(keyToEncrypt, iv...)
	encryptedKey, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		s.publicKey,
		keyToEncrypt,
		nil,
	)
	if err != nil {
		return fmt.Errorf("RSA encryption failed: %w", err)
	}

	var finalPayload []byte
	finalPayload = append(finalPayload, encryptedKey...)
	finalPayload = append(finalPayload, ciphertext...)

	url := s.endpoint

	log.Printf("new request to url=%s, method=%s", url, http.MethodPost)
	for _, m := range mList {
		log.Printf("    data: %s", &m)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(finalPayload))
	if err != nil {
		return fmt.Errorf("sender.go func SendQueryUpdateMetrics(): error create request - %w", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Encoding", "gzip")

	if s.hashKey != nil {
		h := hmac.New(sha256.New, []byte(*s.hashKey))
		h.Write(finalPayload)
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
