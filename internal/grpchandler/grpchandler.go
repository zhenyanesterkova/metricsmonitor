package grpchandler

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/sirupsen/logrus"
	proto "github.com/zhenyanesterkova/metricsmonitor/internal/app/proto/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
)

const (
	ErrServer = "something went wrong"
)

type Repositorie interface {
	// UpdateMetric updates a metric with the given type and value.
	UpdateMetric(metric.Metric) (metric.Metric, error)
	// GetAllMetrics retrieves all available metrics from the storage.
	GetAllMetrics() ([][2]string, error)
	// GetMetricValue retrieves a specific metric from the storage by its name and type.
	GetMetricValue(name, typeMetric string) (metric.Metric, error)
	// Ping checks the availability of the storage.
	Ping() error
	// UpdateManyMetrics updates multiple metrics in the storage.
	UpdateManyMetrics(ctx context.Context, mList []metric.Metric) error
}

type gRPCHandler struct {
	proto.UnimplementedMonitorServer
	repo   Repositorie
	logger logger.LogrusLogger
}

func New(
	rep Repositorie,
	log logger.LogrusLogger,
) *gRPCHandler {
	return &gRPCHandler{
		repo:   rep,
		logger: log,
	}
}

func (h *gRPCHandler) Ping(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	err := h.repo.Ping()
	if err != nil {
		h.logger.LogrusLog.Errorf("failed repo ping: %v", err)
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed ping storage")
	}
	return &emptypb.Empty{}, nil
}

func (h *gRPCHandler) AddMetric(ctx context.Context, req *proto.MetricRequest) (*proto.MetricResponse, error) {
	var newMetric metric.Metric
	switch req.GetMetric().GetType() {
	case metric.TypeCounter:
		newMetric = metric.New(metric.TypeCounter)
		*newMetric.Delta = req.GetMetric().GetDelta()
	case metric.TypeGauge:
		newMetric = metric.New(metric.TypeGauge)
		*newMetric.Value = req.GetMetric().GetValue()
	}
	newMetric.ID = req.GetMetric().GetId()

	h.logger.LogrusLog.WithFields(logrus.Fields{
		"ID":    newMetric.ID,
		"Type":  newMetric.MType,
		"Value": *newMetric.Value,
		"Delta": *newMetric.Delta,
	}).Info("metric for updating in grpc")

	m, err := h.repo.UpdateMetric(newMetric)
	if err != nil {
		switch {
		case errors.Is(err, metric.ErrInvalidName):
			return nil, fmt.Errorf("failed update metric: %w", status.Error(codes.NotFound, "check id of the metric"))
		case errors.Is(err, metric.ErrParseValue) ||
			errors.Is(err, metric.ErrUnknownType) ||
			errors.Is(err, metric.ErrInvalidType):
			return nil,
				fmt.Errorf(
					"failed update metric: %w",
					status.Errorf(codes.InvalidArgument,
						"check value and type of the metric",
					))
		default:
			h.logger.LogrusLog.Errorf("failed update metric: %v", err)
			return nil, fmt.Errorf("failed update metric: %w", status.Errorf(codes.Internal, ErrServer))
		}
	}

	return &proto.MetricResponse{
		Metric: &proto.Metric{
			Id:    m.ID,
			Type:  m.MType,
			Value: *m.Value,
			Delta: *m.Delta,
		},
	}, nil
}
