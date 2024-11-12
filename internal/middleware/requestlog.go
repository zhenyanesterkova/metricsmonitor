package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, fmt.Errorf("logger.go Write() - %w", err)
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func (lm MiddlewareStruct) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := lm.Logger.LogrusLog

		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		defer func() {
			log.WithFields(logrus.Fields{
				"URI":      r.URL.Path,
				"Method":   r.Method,
				"Duration": time.Since(start),
				"Status":   lw.responseData.status,
				"Size":     lw.responseData.size,
			}).Info("got incoming HTTP request")
		}()

		next.ServeHTTP(&lw, r)
	})
}
