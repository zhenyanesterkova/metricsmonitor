package middleware

import (
	"net"
	"net/http"
)

func (lm MiddlewareStruct) CheckTrustIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if lm.trustIPNet != nil {
			ipStr := r.Header.Get("X-Real-IP")

			ip := net.ParseIP(ipStr)

			allowed := lm.trustIPNet.Contains(ip)

			if !allowed {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
