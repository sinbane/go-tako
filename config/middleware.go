package config

import "time"

// CORS middleware config
type CORS struct {
	AllowedOrigins []string `toml:"allowed_origins"`
	AllowedMethods []string `toml:"allowed_methods"`
	AllowedHeaders []string `toml:"allowed_headers"`
}

// CircuitBreaker middleware config
type CircuitBreaker struct {
	MaxRequests  int           `toml:"max_requests"`
	Interval     time.Duration `toml:"interval"`
	Timeout      time.Duration `toml:"timeout"`
	MinRequests  int           `toml:"min_requests"`
	FailureRatio float64       `toml:"failure_ratio"`
}
