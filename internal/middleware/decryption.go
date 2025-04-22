package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"
	"net/http"
)

const (
	cryptoKeySize = 256
	aesKeySize    = 32
)

func isEncryption(cType string) bool {
	return cType == "application/octet-stream"
}

func (lm MiddlewareStruct) DecryptionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPprofPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		contentType := r.Header.Get("Content-Type")
		if !isEncryption(contentType) {
			next.ServeHTTP(w, r)
			return
		}

		encryptedData, err := io.ReadAll(r.Body)
		if err != nil {
			lm.Logger.LogrusLog.Errorf("failed read encrypted data from body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		encryptedKey := encryptedData[:cryptoKeySize]
		ciphertext := encryptedData[cryptoKeySize:]

		decryptedKey, err := rsa.DecryptOAEP(
			sha256.New(),
			rand.Reader,
			lm.privateKey,
			encryptedKey,
			nil,
		)
		if err != nil {
			lm.Logger.LogrusLog.Errorf("failed decrypted key: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		aesKey := decryptedKey[:aesKeySize]
		iv := decryptedKey[aesKeySize:]

		block, err := aes.NewCipher(aesKey)
		if err != nil {
			lm.Logger.LogrusLog.Errorf("failed create AES block: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		stream := cipher.NewCTR(block, iv)
		plaintext := make([]byte, len(ciphertext))
		stream.XORKeyStream(plaintext, ciphertext)

		r.Body = io.NopCloser(bytes.NewBuffer(plaintext))
		next.ServeHTTP(w, r)
	})
}
