package config

import (
	"log"
	"os"
)

type Config struct {
	DSN       string
	JWTSecret string
	Port      string
}

func Load() *Config {
	cfg := &Config{
		DSN:       os.Getenv("DB_DSN"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		Port:      os.Getenv("AUTH_GRPC_PORT"),
	}

	if cfg.DSN == "" || cfg.JWTSecret == "" {
		log.Fatal("missing required environment variables")
	}

	if cfg.Port == "" {
		cfg.Port = "50051"
	}

	return cfg
}
