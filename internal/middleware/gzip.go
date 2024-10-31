package middleware

import (
	"net/http"
	"strings"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/mycompress"
)

func isCompression(cType string) bool {
	if cType == "application/json" ||
		cType == "text/html" {
		return true
	}

	return false
}

func GZipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		supportsGzip := false
		acceptEncoding := r.Header.Values("Accept-Encoding")
		for _, val := range acceptEncoding {
			if strings.Contains(val, "gzip") {
				supportsGzip = true
				break
			}
		}

		contentType := r.Header.Get("Content-Type")
		compressing := isCompression(contentType)

		if supportsGzip && compressing {
			cw := mycompress.NewCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := mycompress.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
