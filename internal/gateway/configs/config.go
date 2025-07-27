package configs

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	AuthGRPCAddr   string
	StatGRPCAddr   string
	LinkGRPCAddr   string
	PublicKey      string
	AllowedOrigins []string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	port := getEnv("GATEWAY_PORT", "8080")
	authAddr := getEnv("AUTH_GRPC_ADDR", "auth:50051")
	statAddr := getEnv("STAT_GRPC_ADDR", "stat:50053")
	linkAddr := getEnv("LINK_GRPC_ADDR", "link:50052")

	pubKeyPath := getEnv("JWT_PUBLIC_KEY_PATH", "/run/secrets/jwt_public.pem")
	key, err := os.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatalf("failed to read public key: %v", err)
	}

	origins := parseCSV(getEnv("GATEWAY_CORS_ORIGINS", "http://localhost:3000"))

	return &Config{
		Port:           port,
		AuthGRPCAddr:   authAddr,
		StatGRPCAddr:   statAddr,
		LinkGRPCAddr:   linkAddr,
		PublicKey:      string(key),
		AllowedOrigins: origins,
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseCSV(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
