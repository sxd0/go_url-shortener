package server

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewGRPCServerWithMiddleware(serverOptions ...grpc.ServerOption) *grpc.Server {
	logger, _ := zap.NewProduction()
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	// opts := []grpc.ServerOption{
	// 	grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
	// 		grpc_zap.UnaryServerInterceptor(logger),
	// 		grpc_recovery.UnaryServerInterceptor(),
	// 	)),
	// 	grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
	// 		grpc_zap.StreamServerInterceptor(logger),
	// 		grpc_recovery.StreamServerInterceptor(),
	// 	)),
	// }

	// 1) Chain-им все ваши Unary interceptors: recovery → logging → JWT → Prometheus
	unary := grpc_middleware.ChainUnaryServer(
		grpc_recovery.UnaryServerInterceptor(),
		grpc_zap.UnaryServerInterceptor(logger),
		grpc_prom.UnaryServerInterceptor,
	)
	// 2) Chain-им все ваши Stream interceptors: recovery → logging → Prometheus
	stream := grpc_middleware.ChainStreamServer(
		grpc_recovery.StreamServerInterceptor(),
		grpc_zap.StreamServerInterceptor(logger),
		grpc_prom.StreamServerInterceptor,
	)

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(unary),
		grpc.StreamInterceptor(stream),
	}

	opts = append(opts, serverOptions...)

	// return grpc.NewServer(opts...)
	grpcServer := grpc.NewServer(opts...)
	grpc_prom.Register(grpcServer)
	return grpcServer
}
