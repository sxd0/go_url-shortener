package handler

import (
	"context"

	"github.com/sxd0/go_url-shortener/internal/link/service"
	"github.com/sxd0/go_url-shortener/proto/gen/go/linkpb"
	"google.golang.org/grpc"
)

type LinkHandler struct {
	linkpb.UnimplementedLinkServiceServer
	service *service.LinkService
}

func RegisterLinkHandler(s *grpc.Server, srv *service.LinkService) {
	linkpb.RegisterLinkServiceServer(s, &LinkHandler{
		service: srv,
	})
}

func (h *LinkHandler) CreateLink(ctx context.Context, req *linkpb.CreateLinkRequest) (*linkpb.LinkResponse, error) {
	return &linkpb.LinkResponse{}, nil
}
