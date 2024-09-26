package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sinbane/tako/config"
	"github.com/sinbane/tako/middleware"
	"github.com/sinbane/tako/route"
)

func TestMainSetup(t *testing.T) {
	// Test config loading
	t.Run("LoadConfig", func(t *testing.T) {
		cfg, err := config.LoadConfig("testdata/config.toml")
		if err != nil {
			t.Fatalf("Failed to load test config: %v", err)
		}
		if cfg.Port != 8080 {
			t.Errorf("Expected port 8080, got %d", cfg.Port)
		}
	})

	// Test router initialization
	t.Run("RouterInitialization", func(t *testing.T) {
		rules := []route.Rule{
			{Prefix: "/test", Target: "http://example.com"},
		}
		router := route.Router(rules)
		if router == nil {
			t.Fatal("Router should not be nil")
		}
	})

	// Test middleware chain
	t.Run("MiddlewareChain", func(t *testing.T) {
		cfg, _ := config.LoadConfig("testdata/config.toml")
		router := route.Router(cfg.Rules)
		handler := middleware.Chain(
			cfg,
			router,
			middleware.Logging,
			middleware.RequestID,
		)

		req, err := http.NewRequest("GET", "/echo/test", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		// Check if the middleware chain is working
		// TODO: This is a basic check, needs to add more specific checks
		if req.Header.Get(middleware.HeaderRequestId) == "" {
			t.Errorf("Expected %s header to be set", middleware.HeaderRequestId)
		}
	})
}

func TestVersion(t *testing.T) {
	if VERSION == "" {
		t.Error("VERSION should not be empty")
	}
	if COMMIT_HASH == "" {
		t.Error("COMMIT_HASH should not be empty")
	}
}
