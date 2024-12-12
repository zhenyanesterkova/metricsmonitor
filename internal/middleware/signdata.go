package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

func (lm MiddlewareStruct) CheckSignData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Header[http.CanonicalHeaderKey("HashSHA256")]; ok {
			signRequestData := r.Header.Get("HashSHA256")
			log := lm.Logger.LogrusLog

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Errorf("middleware: CheckSignData - failed read body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = r.Body.Close()
			if err != nil {
				log.Errorf("middleware: CheckSignData - failed close body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			h := hmac.New(sha256.New, []byte(*lm.hashKey))
			h.Write(bodyBytes)
			sum := h.Sum(nil)

			strSign, err := hex.DecodeString(signRequestData)
			if err != nil {
				log.Errorf("middleware: CheckSignData - failed decode hash: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !hmac.Equal(strSign, sum) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
