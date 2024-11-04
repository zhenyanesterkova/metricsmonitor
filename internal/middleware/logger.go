package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
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
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

type LoggerMiddleware struct {
	Logger logger.LogrusLogger
}

func NewLoggerMiddleware(log logger.LogrusLogger) LoggerMiddleware {
	return LoggerMiddleware{
		Logger: log,
	}
}

func (lm LoggerMiddleware) RequestLogger(next http.Handler) http.Handler {
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
