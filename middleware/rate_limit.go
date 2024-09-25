package middleware

import (
	"net/http"
	"sync"

	"github.com/sinbane/tako/config"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(1, 5) // 1 request per second, with a burst of 5
var ipRateLimiters = make(map[string]*rate.Limiter)
var mu sync.Mutex

func RateLimit(cfg *config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			mu.Lock()
			if _, exists := ipRateLimiters[ip]; !exists {
				ipRateLimiters[ip] = rate.NewLimiter(1, 5)
			}
			limiter := ipRateLimiters[ip]
			mu.Unlock()

			if !limiter.Allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
