package middleware

import (
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/sinbane/tako/config"
)

type Middleware func(http.Handler) http.Handler

type MiddlewareFactory func(*config.Config) Middleware

// Chain applies a series of middlewares to a http.Handler
func Chain(cfg *config.Config, h http.Handler, middlewareFactories ...MiddlewareFactory) http.Handler {
	for i := len(middlewareFactories) - 1; i >= 0; i-- {
		h = WithLogging(middlewareFactories[i](cfg))(h)
	}
	return h
}

// WithLogging wraps a middleware with logging
func WithLogging(next Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			name := getMiddlewareName(next)
			log.Printf("[%s] - handling %s", name, r.URL.Path)
			next(h).ServeHTTP(w, r)
			log.Printf("[%s] - handled %s", name, r.URL.Path)
		})
	}
}

// getMiddlewareName extracts the name of the middleware function
func getMiddlewareName(m interface{}) string {
	// Get the full name of the middleware function, with full package path and .funcN suffix
	name := runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name()
	// Extract the name of the middleware function, without the package path and .funcN suffix
	parts := strings.Split(name, ".")
	if len(parts) > 1 {
		return parts[len(parts)-2] // Return the second-to-last part
	}
	return name // Fallback to full name if splitting fails
}

// A handy function to create a new middleware
func New(h http.Handler) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	}
}
