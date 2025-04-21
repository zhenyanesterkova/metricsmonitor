package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
)

type MiddlewareStruct struct {
	Logger     logger.LogrusLogger
	hashKey    *string
	respData   *responseDataWriter
	privateKey *rsa.PrivateKey
}

func NewMiddlewareStruct(
	log logger.LogrusLogger,
	key *string,
	pathToPrivateKey string,
) (MiddlewareStruct, error) {
	responseData := &responseData{
		status:  0,
		size:    0,
		hashKey: key,
	}

	lw := responseDataWriter{
		responseData: responseData,
	}

	privateKeyPEM, err := os.ReadFile(pathToPrivateKey)
	if err != nil {
		return MiddlewareStruct{}, fmt.Errorf("failed read private key from file: %w", err)
	}

	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return MiddlewareStruct{}, fmt.Errorf("failed parses a private key in PKIX, ASN.1 DER form: %w", err)
	}

	return MiddlewareStruct{
		Logger:     log,
		hashKey:    key,
		respData:   &lw,
		privateKey: privateKey,
	}, nil
}

func (lm MiddlewareStruct) ResetRespDataStruct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPprofPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		lm.respData.responseData.size = 0
		lm.respData.responseData.status = 0
		lm.respData.ResponseWriter = w

		next.ServeHTTP(lm.respData, r)
	})
}

func isPprofPath(path string) bool {
	return strings.Contains(path, "/debug/pprof/")
}
