package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config aggregates all configurable parameters for a service.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
}

type ServerConfig struct {
	GRPCAddr string `yaml:"grpc_addr"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

type JWTConfig struct {
	AccessSecret  string `yaml:"access_secret"`
	RefreshSecret string `yaml:"refresh_secret"`
	AccessTTL     string `yaml:"access_ttl"`  // duration, e.g. "15m"
	RefreshTTL    string `yaml:"refresh_ttl"` // duration, e.g. "720h"
}

// Load reads the YAML file at path and unmarshals it into Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
