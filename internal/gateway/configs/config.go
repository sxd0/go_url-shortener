package configs

import (
	"log"
	"os"
	"strconv"
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
	LogLevel       string

	RedisAddr      string
	RedisEnabled   bool
	RedisDialMs    int
	RedisReadMs    int
	RedisWriteMs   int
	RedisPoolSize  int
	RedisMinIdle   int

	LinkCacheTTL int // seconds

	RLEnabled      bool
	RLTTL          int    // seconds
	RLLimit        int    // max requests per window
	RLKeyMode      string // global | route | route+hash
	TrustedProxies []string

	KafkaAddr           string
	KafkaTopic          string
	KafkaAcks           string
	KafkaBatchSize      int
	KafkaBatchTimeoutMs int
	KafkaCompression    string
	KafkaPublishQueue   int
	KafkaPublishWorkers int
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

	origins := splitCSV(getEnv("GATEWAY_CORS_ORIGINS", "http://localhost:3000"))
	level := getEnv("LOG_LEVEL", "info")

	return &Config{
		Port:           port,
		AuthGRPCAddr:   authAddr,
		StatGRPCAddr:   statAddr,
		LinkGRPCAddr:   linkAddr,
		PublicKey:      string(key),
		AllowedOrigins: origins,
		LogLevel:       level,

		RedisAddr:    getEnv("REDIS_ADDR", "redis:6379"),
		RedisEnabled: mustBool(getEnv("REDIS_ENABLED", "true")),
		RedisDialMs:  mustInt(getEnv("REDIS_DIAL_TIMEOUT_MS", "500")),
		RedisReadMs:  mustInt(getEnv("REDIS_READ_TIMEOUT_MS", "200")),
		RedisWriteMs: mustInt(getEnv("REDIS_WRITE_TIMEOUT_MS", "200")),
		RedisPoolSize: mustInt(getEnv("REDIS_POOL_SIZE", "50")),
		RedisMinIdle:  mustInt(getEnv("REDIS_MIN_IDLE_CONNS", "10")),

		LinkCacheTTL: mustInt(getEnv("LINK_CACHE_TTL", "1800")),

		RLEnabled:      mustBool(getEnv("RL_ENABLED", "true")),
		RLTTL:          mustInt(getEnv("RL_TTL", "60")),
		RLLimit:        mustInt(getEnv("RL_LIMIT", "50")),
		RLKeyMode:      getEnv("RL_KEY_MODE", "route"),
		TrustedProxies: splitCSV(getEnv("TRUSTED_PROXIES", "")),

		KafkaAddr:           getEnv("KAFKA_ADDR", "kafka:9092"),
		KafkaTopic:          getEnv("KAFKA_TOPIC", "link.events"),
		KafkaAcks:           getEnv("KAFKA_ACKS", "1"),
		KafkaBatchSize:      mustInt(getEnv("KAFKA_BATCH_SIZE", "200")),
		KafkaBatchTimeoutMs: mustInt(getEnv("KAFKA_BATCH_TIMEOUT_MS", "100")),
		KafkaCompression:    getEnv("KAFKA_COMPRESSION", "snappy"),
		KafkaPublishQueue:   mustInt(getEnv("KAFKA_PUBLISH_QUEUE", "2048")),
		KafkaPublishWorkers: mustInt(getEnv("KAFKA_PUBLISH_WORKERS", "2")),
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func splitCSV(s string) []string {
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

func mustInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func mustBool(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "1" || s == "yes" || s == "y"
}
