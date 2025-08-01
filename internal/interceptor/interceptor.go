package interceptor

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type InterceptorStruct struct {
	Logger  logger.LogrusLogger
	hashKey *string
}

func NewInterceptorStruct(
	log logger.LogrusLogger,
	key *string,
) (InterceptorStruct, error) {
	return InterceptorStruct{
		Logger:  log,
		hashKey: key,
	}, nil
}

func (is InterceptorStruct) RequestLoggerUnary() grpc.UnaryServerInterceptor {
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

		is.Logger.LogrusLog.WithFields(logrus.Fields{
			"Method":   info.FullMethod,
			"Duration": time.Since(start),
			"Status":   statusCode,
		}).Info("got incoming gRPC request")

		return resp, err
	}
}

func (is InterceptorStruct) CheckSignDataUnary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if is.hashKey == nil {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		hashValues := md.Get("hashsha256")
		if len(hashValues) == 0 {
			return handler(ctx, req)
		}

		signRequestData := hashValues[0]

		data, err := serializeRequest(req)
		if err != nil {
			is.Logger.LogrusLog.Errorf("interceptor: CheckSignData - failed serialize request: %v", err)
			return nil, status.Error(codes.InvalidArgument, "failed to serialize request")
		}

		h := hmac.New(sha256.New, []byte(*is.hashKey))
		h.Write(data)
		sum := h.Sum(nil)

		strSign, err := hex.DecodeString(signRequestData)
		if err != nil {
			is.Logger.LogrusLog.Errorf("interceptor: CheckSignData - failed decode hash: %v", err)
			return nil, status.Error(codes.InvalidArgument, "invalid hash format")
		}

		if !hmac.Equal(strSign, sum) {
			return nil, status.Error(codes.InvalidArgument, "invalid signature")
		}

		return handler(ctx, req)
	}
}

func (is InterceptorStruct) GzipUnary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		hashValues := md.Get("content-encoding")
		if len(hashValues) == 0 {
			return handler(ctx, req)
		}

		for _, v := range hashValues {
			if v == "gzip" {
				err := grpc.SetHeader(ctx, metadata.Pairs("content-encoding", "gzip"))
				if err != nil {
					return nil, fmt.Errorf("failed sets the header metadata to be sent from the server to the client: %w", err)
				}
				err = grpc.SetSendCompressor(ctx, gzip.Name)
				if err != nil {
					return nil, fmt.Errorf("failed sets a compressor for outbound messages from the server: %w", err)
				}
				return handler(ctx, req)
			}
		}

		return handler(ctx, req)
	}
}

func serializeRequest(template any) ([]byte, error) {
	protoMsg, ok := template.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("template is not a proto.Message, got %T", template)
	}

	newMsg, err := proto.Marshal(protoMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal protobuf message: %w", err)
	}

	return newMsg, nil
}
