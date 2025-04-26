package config

import (
	"lowerthirdsapi/internal/helpers"
	"lowerthirdsapi/internal/storage"
)

type Config struct {
	Environment string `envconfig:"ENVIRONMENT"`
	MySQLConfig storage.MySQLConfig
}

func New(envDir string) *Config {
	var cfg Config
	_ = helpers.ProcessConfig(envDir, &cfg)
	return &cfg
}
