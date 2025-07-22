package configs

import (
	"log"
	"os"
)

type Config struct {
	Db   *DbConfig
	Auth AuthConfig
	App  AppConfig
}

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type AuthConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string
}

type AppConfig struct {
	Port string
}

func LoadConfig() *Config {
	db := &DbConfig{
		Host:     getEnv("DB_HOST"),
		Port:     getEnv("DB_PORT"),
		User:     getEnv("DB_USER"),
		Password: getEnv("DB_PASSWORD"),
		Name:     getEnv("DB_NAME"),
		SSLMode:  getEnvWithDefault("DB_SSLMODE", "disable"),
	}

	auth := AuthConfig{
		PrivateKeyPath: getEnv("JWT_PRIVATE_KEY_PATH"),
		PublicKeyPath:  getEnv("JWT_PUBLIC_KEY_PATH"),
	}

	app := AppConfig{
		Port: getEnvWithDefault("AUTH_GRPC_PORT", "50051"),
	}

	return &Config{
		Db:   db,
		Auth: auth,
		App:  app,
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
