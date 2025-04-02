package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func (lm MiddlewareStruct) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPprofPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		defer func() {
			lm.Logger.LogrusLog.WithFields(logrus.Fields{
				"URI":      r.URL.Path,
				"Method":   r.Method,
				"Duration": time.Since(start),
				"Status":   lm.respData.responseData.status,
				"Size":     lm.respData.responseData.size,
			}).Info("got incoming HTTP request")
		}()

		next.ServeHTTP(lm.respData, r)
	})
}
