package interceptor

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InterceptorLogger struct {
	Logger Logger
}

func NewInterceptorLogger(log Logger) *InterceptorLogger {
	return &InterceptorLogger{
		Logger: log,
	}
}

func (l *InterceptorLogger) Intercept() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		statusCode := codes.OK
		if err != nil {
			statusCode = status.Code(err)
		}

		l.Logger.Info("got incoming gRPC request", map[string]any{
			"Method":   info.FullMethod,
			"Duration": time.Since(start),
			"Status":   statusCode,
			"Error":    err,
		})

		return resp, err
	}
}

type LogrusAdapter struct {
	logger logger.LogrusLogger
}

func NewLogrusAdapter(logger logger.LogrusLogger) *LogrusAdapter {
	return &LogrusAdapter{
		logger: logger,
	}
}

func (la *LogrusAdapter) Info(msg string, fields map[string]interface{}) {
	la.logger.LogrusLog.WithFields(logrus.Fields(fields)).Info(msg)
}

func (la *LogrusAdapter) Error(msg string, args ...any) {
	la.logger.LogrusLog.Errorf(msg, args...)
}
