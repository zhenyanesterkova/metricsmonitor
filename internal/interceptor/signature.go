package interceptor

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	opSign = "from signature validation interceptor"
)

type HMACSignatureValidator struct {
	hashKey *string
	logger  Logger
}

func NewHMACSignatureValidator(hashKey *string, logger Logger) *HMACSignatureValidator {
	return &HMACSignatureValidator{
		hashKey: hashKey,
		logger:  logger,
	}
}

func (h *HMACSignatureValidator) ValidateSignature(ctx context.Context, req any) error {
	signRequestData, ok := getMetadataValue(ctx, "hashsha256")
	if !ok {
		return nil
	}

	data, err := serializeRequest(req)
	if err != nil {
		h.logger.Error("%v: failed serialize request: %v", opSign, err)
		return fmt.Errorf("%w", status.Error(codes.InvalidArgument, "failed to serialize request"))
	}

	hmacHash := hmac.New(sha256.New, []byte(*h.hashKey))
	hmacHash.Write(data)
	sum := hmacHash.Sum(nil)

	strSign, err := hex.DecodeString(signRequestData)
	if err != nil {
		h.logger.Error("%v: failed decode hash: %v", opSign, err)
		return fmt.Errorf("%w", status.Error(codes.InvalidArgument, "invalid hash format"))
	}

	if !hmac.Equal(strSign, sum) {
		return fmt.Errorf("%w", status.Error(codes.InvalidArgument, "invalid signature"))
	}

	return nil
}

type SignatureInterceptor struct {
	validator SignatureValidator
}

func NewSignatureInterceptor(validator SignatureValidator) *SignatureInterceptor {
	return &SignatureInterceptor{
		validator: validator,
	}
}

func (s *SignatureInterceptor) Intercept() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if err := s.validator.ValidateSignature(ctx, req); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		return handler(ctx, req)
	}
}
