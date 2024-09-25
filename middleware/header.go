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

func CheckHeaders(cfg *config.Config) Middleware {
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

			next.ServeHTTP(w, r)
		})
	}
}
