package grpchandler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/sirupsen/logrus"
	proto "github.com/zhenyanesterkova/metricsmonitor/internal/app/proto/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
	"github.com/zhenyanesterkova/metricsmonitor/web"
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

func (h *gRPCHandler) AddMetrics(ctx context.Context, req *proto.MetricsRequest) (*emptypb.Empty, error) {
	metricList := make([]metric.Metric, len(req.GetMetrics()))
	for i, m := range req.GetMetrics() {
		var newMetric metric.Metric
		switch m.GetType() {
		case metric.TypeCounter:
			newMetric = metric.New(metric.TypeCounter)
			*newMetric.Delta = m.GetDelta()
		case metric.TypeGauge:
			newMetric = metric.New(metric.TypeGauge)
			*newMetric.Value = m.GetValue()
		}
		newMetric.ID = m.GetId()
		metricList[i] = newMetric
	}

	h.logger.LogrusLog.Info("update metrics in gRPC")

	uCtx := context.WithoutCancel(ctx)
	err := h.repo.UpdateManyMetrics(uCtx, metricList)
	if err != nil {
		h.logger.LogrusLog.Errorf("failed update metrics: %v", err)
		return nil, fmt.Errorf("failed update metrics: %w", status.Errorf(codes.Internal, ErrServer))
	}

	return &emptypb.Empty{}, nil
}

func (h *gRPCHandler) GetMetrics(ctx context.Context, req *emptypb.Empty) (*proto.MetricsHTMLResponse, error) {
	res, err := h.repo.GetAllMetrics()

	if err != nil {
		h.logger.LogrusLog.Errorf("failed get metrics: %v", err)
		return nil, fmt.Errorf("failed get metrics: %w", status.Errorf(codes.Internal, ErrServer))
	}

	tmplMetrics, err := template.ParseFS(web.Templates, "template/allMetricsView.html")
	if err != nil {
		h.logger.LogrusLog.Errorf("failed parse html template: %v", err)
		return nil, fmt.Errorf("failed get metrics: %w", status.Errorf(codes.Internal, ErrServer))
	}

	var buf bytes.Buffer
	err = tmplMetrics.ExecuteTemplate(&buf, "metrics", res)
	if err != nil {
		h.logger.LogrusLog.Errorf("failed execute html template: %v", err)
		return nil, fmt.Errorf("failed get metrics: %w", status.Errorf(codes.Internal, ErrServer))
	}

	return &proto.MetricsHTMLResponse{
		MetricsHTML: buf.String(),
	}, nil
}

func (h *gRPCHandler) GetMetric(ctx context.Context, req *proto.MetricRequest) (*proto.MetricResponse, error) {
	res, err := h.repo.GetMetricValue(req.GetMetric().GetId(), req.GetMetric().GetType())

	if err != nil {
		if errors.Is(err, metric.ErrUnknownMetric) || errors.Is(err, metric.ErrInvalidType) {
			return nil, fmt.Errorf("failed get metric: %w", status.Error(codes.NotFound, "check id and type of the metric"))
		}
		h.logger.LogrusLog.Errorf("get metric value: %v", err)
		return nil, fmt.Errorf("failed get metric: %w", status.Errorf(codes.Internal, ErrServer))
	}

	return &proto.MetricResponse{
		Metric: &proto.Metric{
			Id:    res.ID,
			Type:  res.MType,
			Value: *res.Value,
			Delta: *res.Delta,
		},
	}, nil
}
