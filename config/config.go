package config

import (
	"github.com/BurntSushi/toml"
	"github.com/sinbane/tako/middleware"
	"github.com/sinbane/tako/route"
)

type Config struct {
	Port     int             `toml:"port"`
	Rules    []route.Rule    `toml:"rules"`
	CORS     middleware.CORS `toml:"cors"`
	ServerId string          `toml:"server_id"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
