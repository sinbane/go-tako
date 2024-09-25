package middleware

import (
	"net/http"

	"github.com/sinbane/tako/config"
)

type Middleware func(http.Handler) http.Handler

type MiddlewareFactory func(*config.Config) Middleware

// Chain applies a series of middlewares to a http.Handler
func Chain(cfg *config.Config, h http.Handler, middlewareFactories ...MiddlewareFactory) http.Handler {
	for i := len(middlewareFactories) - 1; i >= 0; i-- {
		h = middlewareFactories[i](cfg)(h)
	}
	return h
}

// New creates a new middleware
func New(h http.Handler) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	}
}
