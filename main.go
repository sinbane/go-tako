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
	router := route.Router()

	// Chain middlewares
	handler := middleware.Chain(
		router,
		middleware.Logging,
		middleware.RateLimit,
		middleware.Cors,
		middleware.RequestID,
		middleware.RemoteIP,
		middleware.CheckHeaders,
	)

	// Start the server
	log.Printf("Starting tako on %d...", cfg.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), handler); err != nil {
		log.Fatalf("Could not start tako: %v", err)
	}
}
