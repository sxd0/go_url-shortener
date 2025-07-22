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
			NewJWTUnaryInterceptor(jwtService),
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(),
		)),
	}

	opts = append(opts, serverOptions...)

	return grpc.NewServer(opts...)
}
