package middleware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/sinbane/tako/config"
)

// count is used to generate a unique request ID for each incoming request
var count uint64

// RequestID is a middleware that generates a unique request ID for each incoming request
func RequestID(cfg *config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			prefix := encode(cfg.ServerId, 28)
			id := atomic.AddUint64(&count, 1)
			r.Header.Set(HeaderRequestId, fmt.Sprintf("%s-%04d", prefix, id))
			next.ServeHTTP(w, r)
		})
	}
}

// Generate a random encoded string of a fixed length
func encode(input string, fixedLength int) string {
	var buf [12]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		log.Printf("Error generating random bytes: %v\n", err)
	}
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
