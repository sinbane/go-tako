package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/sinbane/tako/config"
)

func Logging(cfg *config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Printf("%s %s took %v", r.Method, r.URL.Path, time.Since(start))
		})
	}
}
