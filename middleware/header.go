package middleware

import (
	"fmt"
	"net/http"

	"github.com/sinbane/tako/config"
)

const (
	HeaderApiVersion = "Tako-API-Version"
	HeaderCustomerId = "Tako-Customer-Id"
	HeaderRequestId  = "Tako-Request-Id"
	HeaderUserAgent  = "Tako-User-Agent"
)

var protectedHeaders = []string{
	HeaderCustomerId,
	HeaderRequestId,
}

var requiredHeaders = []string{
	HeaderApiVersion,
	"Authorization",
}

func Header(cfg *config.Config) Middleware {
	// add to protected headers from config
	for _, header := range cfg.Header.Protected {
		protectedHeaders = append(protectedHeaders, header)
	}

	// add to required headers from config
	for _, header := range cfg.Header.Required {
		requiredHeaders = append(requiredHeaders, header)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, header := range protectedHeaders {
				if r.Header.Get(header) != "" {
					http.Error(w, fmt.Sprintf("Header %s is protected", header), http.StatusBadRequest)
					return
				}
			}

			for _, header := range requiredHeaders {
				if r.Header.Get(header) == "" {
					http.Error(w, fmt.Sprintf("Header %s is required", header), http.StatusBadRequest)
					return
				}
			}

			// Check if the request URL is in the bypass URLs
			for _, bypassURL := range cfg.Header.BypassURLs {
				if r.URL.Path == bypassURL {
					next.ServeHTTP(w, r)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
