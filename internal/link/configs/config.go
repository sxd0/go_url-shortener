package configs

import (
	"log"
	"os"
)

type Config struct {
	Db  *DbConfig
	App AppConfig
}

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type AppConfig struct {
	Port string
}

func LoadConfig() *Config {
	return &Config{
		Db: &DbConfig{
			Host:     getEnv("DB_HOST"),
			Port:     getEnv("DB_PORT"),
			User:     getEnv("DB_USER"),
			Password: getEnv("DB_PASSWORD"),
			Name:     getEnv("DB_NAME"),
			SSLMode:  getEnvWithDefault("DB_SSLMODE", "disable"),
		},
		App: AppConfig{
			Port: getEnvWithDefault("LINK_GRPC_PORT", "50052"),
		},
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s not set", key)
	}
	return val
}

func getEnvWithDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
