package server

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sxd0/go_url-shortener/internal/auth/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewGRPCServerWithMiddleware(jwtService *jwt.JWT, serverOptions ...grpc.ServerOption) *grpc.Server {
	logger, _ := zap.NewProduction()
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// ВАЖНО: recovery должен быть СНАЧАЛА, чтобы ловить паники из всех нижележащих интерсепторов
			grpc_recovery.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger),
			NewJWTUnaryInterceptor(jwtService),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger),
		)),
	}

	opts = append(opts, serverOptions...)
	return grpc.NewServer(opts...)
}
