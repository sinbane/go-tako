package middleware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
)

// RequestID is a middleware that generates a unique request ID for each incoming request
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := encode(r.RemoteAddr, 16)
		w.Header().Set(HeaderRequestId, requestID)
		next.ServeHTTP(w, r)
	})
}

// Generate a encoded string of a fixed length
func encode(input string, fixedLength int) string {
	var buf [12]byte
	rand.Read(buf[:])
	hasher := sha256.New()
	hasher.Write(append([]byte(input), buf[:]...))
	hashBytes := hasher.Sum(nil)

	b64 := base64.StdEncoding.EncodeToString(hashBytes)
	b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)

	if len(b64) > fixedLength {
		b64 = b64[:fixedLength] // Trim if longer
	} else {
		for len(b64) < fixedLength {
			b64 += "="
		}
	}

	return b64
}
