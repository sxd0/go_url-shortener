package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	AuthGRPCAddr string
	StatGRPCAddr string
	LinkGRPCAddr string
	PublicKey    string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	authAddr := os.Getenv("AUTH_GRPC_ADDR")
	statAddr := os.Getenv("STAT_GRPC_ADDR")
	linkAddr := os.Getenv("LINK_GRPC_ADDR")
	pubKeyPath := os.Getenv("JWT_PUBLIC_KEY_PATH")

	if pubKeyPath == "" {
		log.Fatal("JWT_PUBLIC_KEY_PATH is required")
	}

	key, err := os.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatalf("failed to read public key: %v", err)
	}

	return &Config{
		Port:         port,
		AuthGRPCAddr: authAddr,
		StatGRPCAddr: statAddr,
		LinkGRPCAddr: linkAddr,
		PublicKey:    string(key),
	}
}
