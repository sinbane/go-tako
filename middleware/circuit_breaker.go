package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sinbane/tako/config"
	"github.com/sony/gobreaker"
)

type PanelBoard struct {
	breakers map[string]*gobreaker.CircuitBreaker
}

func NewPanelBoard(cfg *config.Config) *PanelBoard {
	breakers := make(map[string]*gobreaker.CircuitBreaker)

	for id, _cb := range cfg.CircuitBreakers {
		breakers[id] = gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        id,
			MaxRequests: uint32(_cb.MaxRequests),
			Interval:    time.Duration(_cb.Interval) * time.Second,
			Timeout:     time.Duration(_cb.Timeout) * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.Requests >= uint32(_cb.MinRequests) && failureRatio >= _cb.FailureRatio
			},
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				fmt.Printf("Circuit Breaker '%s' state changed from %s to %s\n", name, from, to)
			},
		})
	}

	return &PanelBoard{breakers: breakers}
}

func CircuitBreaker(cfg *config.Config) Middleware {
	board := NewPanelBoard(cfg)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Determine which circuit breaker to use based on the request
			breakerID := "default"
			for _, rule := range cfg.Rules {
				if rule.Prefix != "" && strings.HasPrefix(r.URL.Path, rule.Prefix) {
					if rule.ID != "" {
						breakerID = rule.ID
					}
					break
				}
			}
			breaker, exists := board.breakers[breakerID]
			if !exists {
				next.ServeHTTP(w, r)
				return
			}

			_, err := breaker.Execute(func() (interface{}, error) {
				rw := &responseWriter{ResponseWriter: w}
				next.ServeHTTP(rw, r)
				if rw.status >= 500 {
					return nil, fmt.Errorf("server error: %d", rw.status)
				}
				return nil, nil
			})

			if err != nil {
				if w.Header().Get("Content-Type") == "" {
					// No error has been written yet
					http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
				}
				return
			}
		})
	}
}

// responseWriter is a custom ResponseWriter that captures the status code since
// the original ResponseWriter has no method to get the status code in gobreaker.Execute()
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
