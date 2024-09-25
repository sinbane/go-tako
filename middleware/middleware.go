package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

type Middlewares []func(http.Handler) http.Handler

// Chain applies a series of middlewares to a http.Handler
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
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
