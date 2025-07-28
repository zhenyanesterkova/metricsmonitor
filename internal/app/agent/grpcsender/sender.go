package grpcsender

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
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
	pb "github.com/zhenyanesterkova/metricsmonitor/internal/app/proto/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	aesKeySize = 32
	op         = "grpc sender: "
)

type GRPCSender struct {
	report                  ReportData
	client                  pb.MonitorClient
	conn                    *grpc.ClientConn
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
) (*GRPCSender, error) {
	var publicKeyRsa *rsa.PublicKey
	if pathToPublicKey != "" {
		publicKeyPEM, err := os.ReadFile(pathToPublicKey)
		if err != nil {
			return nil, fmt.Errorf("%s failed read public key from file: %w", op, err)
		}

		publicKeyBlock, _ := pem.Decode(publicKeyPEM)
		publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("%s failed parses a public key in PKIX, ASN.1 DER form: %w", op, err)
		}

		var ok bool
		publicKeyRsa, ok = publicKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("%s %w", op, errors.New("failed converting type to *rsa.PublicKey"))
		}
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s failed to connect to gRPC server: %w", op, err)
	}

	client := pb.NewMonitorClient(conn)

	return &GRPCSender{
		client:         client,
		conn:           conn,
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

func (s *GRPCSender) Close() error {
	err := s.conn.Close()
	if err != nil {
		return fmt.Errorf("%s failed to close gRPC conn: %w", op, err)
	}

	return nil
}

func (s *GRPCSender) SendQueryUpdateMetrics() error {
	pbMetrics := s.report.metricsBuf.GetMetricsListForGRPC()

	if len(pbMetrics) == 0 {
		log.Printf("    %s no data for grpc sending ...\n", op)
		return nil
	}

	req := &pb.MetricsRequest{
		Metrics: pbMetrics,
	}

	ctx := context.Background()
	if s.publicKey != nil || s.hashKey != nil {
		jsonData, err := json.Marshal(pbMetrics)
		if err != nil {
			return fmt.Errorf("%s failed to marshal metrics to JSON: %w", op, err)
		}

		var buff bytes.Buffer
		gzWriter := gzip.NewWriter(&buff)
		if _, err := gzWriter.Write(jsonData); err != nil {
			return fmt.Errorf("%s failed to compress data: %w", op, err)
		}
		if err := gzWriter.Close(); err != nil {
			return fmt.Errorf("%s failed to close gzip writer: %w", op, err)
		}

		compressedData := buff.Bytes()
		finalData := compressedData

		if s.publicKey != nil {
			aesKey := make([]byte, aesKeySize)
			if _, err := rand.Read(aesKey); err != nil {
				return fmt.Errorf("%s failed to generate AES key: %w", op, err)
			}

			iv := make([]byte, aes.BlockSize)
			if _, err := rand.Read(iv); err != nil {
				return fmt.Errorf("%s failed to generate IV: %w", op, err)
			}

			block, _ := aes.NewCipher(aesKey)
			stream := cipher.NewCTR(block, iv)

			ciphertext := make([]byte, len(compressedData))
			stream.XORKeyStream(ciphertext, compressedData)

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
				return fmt.Errorf("%s RSA encryption failed: %w", op, err)
			}

			finalData = append(finalData, encryptedKey...)
			finalData = append(finalData, ciphertext...)
		}

		md := metadata.New(map[string]string{
			"content-encoding": "gzip",
		})

		if s.hashKey != nil {
			h := hmac.New(sha256.New, []byte(*s.hashKey))
			h.Write(finalData)
			sum := hex.EncodeToString(h.Sum(nil))
			md.Set("hash-sha256", sum)
		}

		if s.publicKey != nil {
			md.Set("encrypted", "true")
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	log.Printf("%s new gRPC request to %s", op, s.endpoint)
	for _, m := range pbMetrics {
		log.Printf("    %s data: %+v", op, m)
	}

	log.Printf("%s send gRPC request ...\n", op)
	_, err := s.client.AddMetrics(ctx, req)

	if err != nil {
		reqSuccess := false
		for i, interval := range s.requestAttemptIntervals {
			dur, errParse := time.ParseDuration(interval)
			if errParse != nil {
				return fmt.Errorf(`%s failed send statistic to server: %w;
				the attempt to re-send № %d failed: 
				the interval could not be parsed: %w`,
					op,
					err,
					i+1,
					errParse,
				)
			}
			time.Sleep(dur)
			_, err = s.client.AddMetrics(ctx, req)
			if err == nil {
				reqSuccess = true
				break
			}
		}
		if !reqSuccess {
			return fmt.Errorf(`%s failed send statistic to server: %w,
			all attempts to re-send failed`,
				op,
				err,
			)
		}
	}

	return nil
}

func (s *GRPCSender) SendReport(ctx context.Context) {
	jobs := make(chan struct{}, 1)

	defer close(jobs)

	for w := 1; w <= s.rateLimit; w++ {
		go s.sendWorker(jobs)
	}

	ticker := time.NewTicker(s.reportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stop grpc send workers.")
			return
		case <-ticker.C:
			select {
			case jobs <- struct{}{}:
				log.Println("Start grpc send statistic ...")
			default:
				log.Println("All grpc workers are busy, skipping this tick")
			}
		}
	}
}

func (s *GRPCSender) sendWorker(jobs <-chan struct{}) {
	for range jobs {
		err := s.SendQueryUpdateMetrics()
		if err != nil {
			log.Printf("%s error occurred while sending metrics to server %v", op, err)
			continue
		}
		s.report.metricsBuf.ResetCountersValues()
	}
}
