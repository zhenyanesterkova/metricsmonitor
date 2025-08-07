package grpcsender

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/retry"
	pb "github.com/zhenyanesterkova/metricsmonitor/internal/app/proto/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
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
	creds, err := credentials.NewClientTLSFromFile(pathToPublicKey, "")
	if err != nil {
		return nil, fmt.Errorf("%s failed to constructs TLS credentials from the provided root certificate: %w", op, err)
	}

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(
			creds,
		),
		grpc.WithDefaultCallOptions(
			grpc.UseCompressor(gzip.Name),
		),
	)
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

	md := metadata.New(map[string]string{
		"content-encoding": "gzip",
	})

	if s.hashKey != nil {
		data, err := proto.Marshal(req)
		if err != nil {
			return fmt.Errorf("%s failed to marshal metrics: %w", op, err)
		}
		h := hmac.New(sha256.New, []byte(*s.hashKey))
		h.Write(data)
		sum := hex.EncodeToString(h.Sum(nil))
		md.Set("hashsha256", sum)
	}

	md.Set("encrypted", "true")

	ctx = metadata.NewOutgoingContext(ctx, md)

	log.Printf("%s new gRPC request to %s", op, s.endpoint)
	for _, m := range pbMetrics {
		log.Printf("    %s data: %+v", op, m)
	}

	log.Printf("%s send gRPC request ...\n", op)
	err := retry.RetryRequest(func() error {
		_, err := s.client.AddMetrics(ctx, req)
		return fmt.Errorf("failed add metrics to server: %w", err)
	}, s.requestAttemptIntervals)

	return fmt.Errorf("failed retry add metrics to server: %w", err)
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
