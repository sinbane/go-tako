package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sinbane/tako/config"
	"github.com/sinbane/tako/middleware"
	"github.com/sinbane/tako/route"
)

func main() {
	log.Printf("Version: %s, Commit ID: %s", VERSION, COMMIT_HASH)

	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize the router
	router := route.Router(cfg.Rules)

	// Chain middlewares
	handler := middleware.Chain(
		cfg,
		router, //the last handler in the chain
		middleware.Logging,
		middleware.CheckHeaders,
		middleware.RateLimit,
		middleware.Cors,
		middleware.RequestID,
		middleware.RemoteIP,
		middleware.CircuitBreaker,
		middleware.JWT,
	)

	// Start the server
	log.Printf("Starting tako on %d...", cfg.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), handler); err != nil {
		log.Fatalf("Could not start tako: %v", err)
	}
}
