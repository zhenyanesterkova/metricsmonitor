package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
)

type responseData struct {
	hashKey *string
	status  int
	size    int
}

type responseDataWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *responseDataWriter) Write(b []byte) (int, error) {
	if r.responseData.hashKey != nil {
		h := hmac.New(sha256.New, []byte(*r.responseData.hashKey))
		h.Write(b)
		sum := hex.EncodeToString(h.Sum(nil))
		r.ResponseWriter.Header().Set("HashSHA256", sum)
	}

	size, err := r.ResponseWriter.Write(b)
	if err != nil {
		return size, fmt.Errorf("logger.go Write() - %w", err)
	}
	r.responseData.size += size

	return size, nil
}

func (r *responseDataWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
