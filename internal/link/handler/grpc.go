package handler

import (
	"context"

	"github.com/sxd0/go_url-shortener/internal/link/model"
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
	link, err := h.service.CreateLink(ctx, req.Url, uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &linkpb.LinkResponse{
		Link: toProto(link),
	}, nil
}

func (h *LinkHandler) GetAllLinks(ctx context.Context, req *linkpb.GetAllLinksRequest) (*linkpb.GetAllLinksResponse, error) {
	links, count, err := h.service.GetAllLinks(ctx, uint(req.UserId), int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	var protoLinks []*linkpb.Link
	for _, l := range links {
		protoLinks = append(protoLinks, toProto(&l))
	}

	return &linkpb.GetAllLinksResponse{
		Links: protoLinks,
		Total: uint64(count),
	}, nil
}

func (h *LinkHandler) UpdateLink(ctx context.Context, req *linkpb.UpdateLinkRequest) (*linkpb.LinkResponse, error) {
	link, err := h.service.UpdateLink(ctx, uint(req.Id), req.Url, req.Hash)
	if err != nil {
		return nil, err
	}
	return &linkpb.LinkResponse{
		Link: toProto(link),
	}, nil
}

func (h *LinkHandler) DeleteLink(ctx context.Context, req *linkpb.DeleteLinkRequest) (*linkpb.Empty, error) {
	if err := h.service.DeleteLink(ctx, uint(req.Id)); err != nil {
		return nil, err
	}
	return &linkpb.Empty{}, nil
}

func (h *LinkHandler) GetLinkByHash(ctx context.Context, req *linkpb.GetLinkByHashRequest) (*linkpb.LinkResponse, error) {
	link, err := h.service.GetLinkByHash(ctx, req.Hash)
	if err != nil {
		return nil, err
	}
	return &linkpb.LinkResponse{
		Link: toProto(link),
	}, nil
}

func toProto(link *model.Link) *linkpb.Link {
	return &linkpb.Link{
		Id:     uint32(link.ID),
		Url:    link.Url,
		Hash:   link.Hash,
		UserId: uint32(link.UserID),
	}
}
