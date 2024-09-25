package route

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

func ReverseProxy(rule Rule) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		originalURL := r.URL
		targetURL, err := url.Parse(rule.Target)
		if err != nil {
			log.Printf("Error parsing target url: %v", err)
			http.Error(w, "Bad gateway", http.StatusBadGateway)
			return
		}

		// Preserve the suffix of the original path
		targetURL.Path = path.Join(targetURL.Path, strings.TrimPrefix(originalURL.Path, rule.Prefix))
		targetURL.RawQuery = originalURL.RawQuery // Preserve query parameters

		// Create a new request to the target URL
		req, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
		if err != nil {
			log.Printf("Error creating new request: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Copy all headers from the original request
		req.Header = r.Header

		// Perform the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Error performing request: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Copy the response headers and status code
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)

		// Copy the response body
		io.Copy(w, resp.Body)
	}
}
