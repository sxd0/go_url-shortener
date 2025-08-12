package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sxd0/go_url-shortener/internal/gateway"
	"github.com/sxd0/go_url-shortener/internal/gateway/configs"
	"github.com/sxd0/go_url-shortener/internal/gateway/jwt"
	"github.com/sxd0/go_url-shortener/internal/gateway/logger"
	"github.com/sxd0/go_url-shortener/internal/gateway/redis"
	"github.com/sxd0/go_url-shortener/internal/gateway/service"
	"github.com/sxd0/go_url-shortener/pkg/kafka"
)

func main() {
	cfg := configs.LoadConfig()

	logger.InitFromEnv()
	defer logger.Sync()

	verifier := jwt.NewVerifier(cfg.PublicKey)

	authSvc, err := service.NewAuthService(cfg.AuthGRPCAddr)
	if err != nil {
		log.Fatalf("auth client: %v", err)
	}
	linkSvc, err := service.NewLinkService(cfg.LinkGRPCAddr)
	if err != nil {
		log.Fatalf("link client: %v", err)
	}
	statSvc, err := service.NewStatService(cfg.StatGRPCAddr)
	if err != nil {
		log.Fatalf("stat client: %v", err)
	}

	producer := kafka.NewPublisherWithConfig(kafka.PubConfig{
		Brokers:      []string{cfg.KafkaAddr},
		Topic:        cfg.KafkaTopic,
		Acks:         cfg.KafkaAcks,                  // "1" по умолчанию
		BatchSize:    cfg.KafkaBatchSize,             // 200
		BatchTimeout: time.Duration(cfg.KafkaBatchTimeoutMs) * time.Millisecond,
		Compression:  cfg.KafkaCompression,           // snappy
		QueueSize:    cfg.KafkaPublishQueue,          // 2048
		Workers:      cfg.KafkaPublishWorkers,        // 2
	})
	defer producer.Close()

	var rdb *redis.Client
	if cfg.RedisEnabled {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.RedisDialMs)*time.Millisecond)
		defer cancel()
		rdb, err = redis.New(ctx, redis.Options{
			Addr:         cfg.RedisAddr,
			DialTimeout:  time.Duration(cfg.RedisDialMs) * time.Millisecond,
			ReadTimeout:  time.Duration(cfg.RedisReadMs) * time.Millisecond,
			WriteTimeout: time.Duration(cfg.RedisWriteMs) * time.Millisecond,
			PoolSize:     cfg.RedisPoolSize,
			MinIdleConns: cfg.RedisMinIdle,
		})
		if err != nil {
			log.Printf("WARNING: redis disabled (connect failed): %v", err)
			rdb = nil
		}
	} else {
		log.Printf("INFO: redis disabled by config")
	}

	router := gateway.NewRouter(gateway.RedirectDeps{
		Verifier:   verifier,
		AuthClient: authSvc.Client(),
		LinkClient: linkSvc.Client(),
		StatClient: statSvc.Client(),
		Cache:      rdb,
		KafkaPublisher: producer,
	}, cfg)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Gateway listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}
	if rdb != nil {
		_ = rdb.Close()
	}
}
