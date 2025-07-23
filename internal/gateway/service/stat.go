package service

import (
	"context"
	"fmt"

	"github.com/sxd0/go_url-shortener/proto/gen/go/statpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type StatService struct {
	client statpb.StatServiceClient
}

func NewStatService(addr string) (*StatService, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to StatService: %w", err)
	}

	client := statpb.NewStatServiceClient(conn)
	return &StatService{client: client}, nil
}

func (s *StatService) GetStats(ctx context.Context, userID uint, from, to, by string) (*statpb.GetStatsResponseList, error) {
	md := metadata.Pairs("x-user-id", fmt.Sprintf("%d", userID))
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &statpb.GetStatsRequest{
		From: from,
		To:   to,
		By:   by,
	}

	return s.client.GetStats(ctx, req)
}
