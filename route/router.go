package route

import (
	"net/http"
)

func Router(rules []Rule) http.Handler {
	mux := http.NewServeMux()

	// Proxy to backend microservices based on rules
	for _, rule := range rules {
		mux.Handle(rule.Prefix, ReverseProxy(rule.Target))
	}

	return mux
}

type Rule struct {
	Prefix string
	Target string
}
