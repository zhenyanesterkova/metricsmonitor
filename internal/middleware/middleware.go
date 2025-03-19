package middleware

import (
	"net/http"
	"strings"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
)

type MiddlewareStruct struct {
	Logger   logger.LogrusLogger
	hashKey  *string
	respData *responseDataWriter
}

func NewMiddlewareStruct(log logger.LogrusLogger, key *string) MiddlewareStruct {
	responseData := &responseData{
		status:  0,
		size:    0,
		hashKey: key,
	}

	lw := responseDataWriter{
		responseData: responseData,
	}

	return MiddlewareStruct{
		Logger:   log,
		hashKey:  key,
		respData: &lw,
	}
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
