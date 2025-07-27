package handler

import (
	"context"
	"errors"

	"github.com/sxd0/go_url-shortener/internal/link/model"
	"github.com/sxd0/go_url-shortener/internal/link/service"
	"github.com/sxd0/go_url-shortener/proto/gen/go/linkpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type LinkHandler struct {
	linkpb.UnimplementedLinkServiceServer
	service *service.LinkService
}

func NewLinkHandler(s *service.LinkService) *LinkHandler {
	return &LinkHandler{service: s}
}

func (h *LinkHandler) CreateLink(ctx context.Context, req *linkpb.CreateLinkRequest) (*linkpb.LinkResponse, error) {
	if req.GetUrl() == "" || req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "url and user_id are required")
	}
	link, err := h.service.CreateLink(ctx, req.Url, uint(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create link")
	}
	return &linkpb.LinkResponse{Link: toProto(link)}, nil
}

func (h *LinkHandler) GetAllLinks(ctx context.Context, req *linkpb.GetAllLinksRequest) (*linkpb.GetAllLinksResponse, error) {
	links, total, err := h.service.GetAllLinks(ctx, uint(req.UserId), int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get links")
	}

	out := make([]*linkpb.Link, 0, len(links))
	for _, l := range links {
		out = append(out, toProto(&l))
	}
	return &linkpb.GetAllLinksResponse{
		Links: out,
		Total: uint64(total),
	}, nil
}

func (h *LinkHandler) UpdateLink(ctx context.Context, req *linkpb.UpdateLinkRequest) (*linkpb.LinkResponse, error) {
	if req.GetUrl() == "" && req.GetHash() == "" && req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id or hash and url are required")
	}

	link, err := h.service.UpdateLink(ctx, uint(req.Id), req.Url, req.Hash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "link not found" {
			return nil, status.Error(codes.NotFound, "link not found")
		}
		return nil, status.Error(codes.Internal, "failed to update link")
	}
	return &linkpb.LinkResponse{Link: toProto(link)}, nil
}

func (h *LinkHandler) DeleteLink(ctx context.Context, req *linkpb.DeleteLinkRequest) (*linkpb.Empty, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if err := h.service.DeleteLink(ctx, uint(req.Id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "link not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete link")
	}
	return &linkpb.Empty{}, nil
}

func (h *LinkHandler) GetLinkByHash(ctx context.Context, req *linkpb.GetLinkByHashRequest) (*linkpb.LinkResponse, error) {
	if req.GetHash() == "" {
		return nil, status.Error(codes.InvalidArgument, "hash is required")
	}
	link, err := h.service.GetLinkByHash(ctx, req.Hash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "link not found")
		}
		return nil, status.Error(codes.Internal, "failed to get link")
	}
	return &linkpb.LinkResponse{Link: toProto(link)}, nil
}

func RegisterGRPC(server *grpc.Server, h *LinkHandler) {
	linkpb.RegisterLinkServiceServer(server, h)
}

func RegisterLinkHandler(server *grpc.Server, h *LinkHandler) {
	linkpb.RegisterLinkServiceServer(server, h)
}

func toProto(link *model.Link) *linkpb.Link {
	return &linkpb.Link{
		Id:     uint32(link.ID),
		Url:    link.Url,
		Hash:   link.Hash,
		UserId: uint32(link.UserID),
	}
}
