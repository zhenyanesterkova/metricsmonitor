package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

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

func getMetadataValue(ctx context.Context, key string) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	values := md.Get(key)
	if len(values) == 0 {
		return "", false
	}

	return values[0], true
}

func getMetadataValues(ctx context.Context, key string) ([]string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return []string{}, false
	}

	values := md.Get(key)
	if len(values) == 0 {
		return []string{}, false
	}

	return values, true
}
