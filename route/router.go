package route

import (
	"net/http"
)

func Router() http.Handler {
	mux := http.NewServeMux()

	// Proxy to backend microservices
	mux.Handle("/api/v1/users/", ReverseProxy("http://localhost:8081"))
	mux.Handle("/api/v1/orders/", ReverseProxy("http://localhost:8082"))

	return mux
}

type Rule struct {
	Prefix string
	Target string
}
