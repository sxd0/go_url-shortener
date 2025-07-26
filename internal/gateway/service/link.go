package service

import (
	"fmt"

	"github.com/sxd0/go_url-shortener/proto/gen/go/linkpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LinkService struct {
	client linkpb.LinkServiceClient
}

func (s *LinkService) Client() linkpb.LinkServiceClient {
	return s.client
}

func NewLinkService(addr string) (*LinkService, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to LinkService: %w", err)
	}

	client := linkpb.NewLinkServiceClient(conn)
	return &LinkService{client: client}, nil
}
