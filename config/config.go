package config

import (
	"github.com/BurntSushi/toml"
	"github.com/sinbane/tako/route"
)

type Config struct {
	Port     int          `toml:"port"`
	Rules    []route.Rule `toml:"rules"`
	CORS     CORS         `toml:"cors"`
	ServerId string       `toml:"server_id"`
}

type CORS struct {
	AllowedOrigins []string `toml:"allowed_origins"`
	AllowedMethods []string `toml:"allowed_methods"`
	AllowedHeaders []string `toml:"allowed_headers"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
