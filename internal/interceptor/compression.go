package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	opCompress      = "from compression interceptor"
	metadataEncName = "content-encoding"
	errInternal     = "something went wrong"
)

type GzipCompressionHandler struct {
	logger Logger
}

func NewGzipCompressionHandler(logger Logger) *GzipCompressionHandler {
	return &GzipCompressionHandler{
		logger: logger,
	}
}

func (g *GzipCompressionHandler) HandleCompression(ctx context.Context) error {
	hashValues, ok := getMetadataValues(ctx, metadataEncName)
	if !ok {
		return nil
	}

	for _, v := range hashValues {
		if v == gzip.Name {
			err := grpc.SetHeader(ctx, metadata.Pairs(metadataEncName, gzip.Name))
			if err != nil {
				g.logger.Error("%v: failed set header %v: %v", opCompress, metadataEncName, err)
				return fmt.Errorf("%w", status.Error(codes.Internal, errInternal))
			}
			err = grpc.SetSendCompressor(ctx, gzip.Name)
			if err != nil {
				g.logger.Error("%v: failed sets a compressor for outbound messages from the server: %v", opCompress, err)
				return fmt.Errorf("%w", status.Error(codes.Internal, errInternal))
			}
			break
		}
	}

	return nil
}

type CompressionInterceptor struct {
	handler CompressionHandler
}

func NewCompressionInterceptor(handler CompressionHandler) *CompressionInterceptor {
	return &CompressionInterceptor{
		handler: handler,
	}
}

func (c *CompressionInterceptor) Intercept() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if err := c.handler.HandleCompression(ctx); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		return handler(ctx, req)
	}
}
