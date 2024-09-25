package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/sinbane/tako/config"
)

func RemoteIP(cfg *config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ip := getRemoteIP(r); ip != "" {
				r.RemoteAddr = ip
			}
			next.ServeHTTP(w, r)
		})
	}
}

var ipHeaders = []string{
	"True-Client-IP", // Cloudflare Enterprise plan
	"X-Real-IP",
	"X-Forwarded-For",
}

func getRemoteIP(r *http.Request) string {
	for _, header := range ipHeaders {
		if ip := r.Header.Get(header); ip != "" {
			ips := strings.Split(ip, ",")
			if ips[0] == "" || net.ParseIP(ips[0]) == nil {
				continue
			}
			return ips[0]
		}
	}
	return ""
}
