package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/sinbane/tako/route"
)

type Config struct {
	Port            int                       `toml:"port"`
	Rules           []route.Rule              `toml:"rules"`
	CORS            CORS                      `toml:"cors"`
	ServerId        string                    `toml:"server_id"`
	CircuitBreakers map[string]CircuitBreaker `toml:"circuit_breakers"`
	JWT             JWT                       `toml:"jwt"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, err
	}

	serverId, ok := os.LookupEnv("K8S_POD_NAME")
	if ok {
		config.ServerId = serverId
	}

	return &config, nil
}
