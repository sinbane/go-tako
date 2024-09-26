package middleware

import (
	"net/http"
	"sync"

	"github.com/sinbane/tako/config"
	"golang.org/x/time/rate"
)

// limiters is a map of IP addresses to their respective rate limiters.
var limiters = make(map[string]*rate.Limiter)
var mu sync.Mutex

func RateLimit(cfg *config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			mu.Lock()
			if _, exists := limiters[ip]; !exists {
				limiters[ip] = rate.NewLimiter(1, 5) // 1 request per second, with a burst of 5
			}
			mu.Unlock()

			limiter := limiters[ip]

			if !limiter.Allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
