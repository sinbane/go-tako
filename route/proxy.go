package route

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func ReverseProxy(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetURL, err := url.Parse(target)
		if err != nil {
			http.Error(w, "Bad gateway", http.StatusBadGateway)
			return
		}

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
