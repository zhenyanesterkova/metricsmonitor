package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

type Logger interface {
	Info(msg string, fields map[string]any)
	Error(msg string, args ...any)
}

type SignatureValidator interface {
	ValidateSignature(ctx context.Context, req any) error
}

type CompressionHandler interface {
	HandleCompression(ctx context.Context) error
}

type UnaryInterceptor interface {
	Intercept() grpc.UnaryServerInterceptor
}
