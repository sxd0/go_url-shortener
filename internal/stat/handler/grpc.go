package handler

import (
	"context"
	"time"

	"strconv"

	"github.com/sxd0/go_url-shortener/internal/stat/repository"
	"github.com/sxd0/go_url-shortener/proto/gen/go/statpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	statpb.UnimplementedStatServiceServer
	Repo *repository.StatRepository
}

func NewStatGRPCHandler(repo *repository.StatRepository) *GRPCHandler {
	return &GRPCHandler{
		Repo: repo,
	}
}

func (h *GRPCHandler) GetStats(ctx context.Context, req *statpb.GetStatsRequest) (*statpb.GetStatsResponseList, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	userIDs := md.Get("x-user-id")
	if len(userIDs) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing x-user-id header")
	}

	userID, err := strconv.ParseUint(userIDs[0], 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid x-user-id")
	}

	from, err := time.Parse("2006-01-02", req.GetFrom())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid 'from' date")
	}

	to, err := time.Parse("2006-01-02", req.GetTo())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid 'to' date")
	}

	by := req.GetBy()
	if by != repository.GroupByDay && by != repository.GroupByMonth {
		return nil, status.Error(codes.InvalidArgument, "invalid 'by' value (expected: day or month)")
	}

	stats := h.Repo.GetStats(uint(userID), by, from, to)

	var resp statpb.GetStatsResponseList
	for _, stat := range stats {
		resp.Stats = append(resp.Stats, &statpb.Stat{
			Period: stat.Period,
			Sum:    int32(stat.Sum),
		})
	}

	return &resp, nil
}

func (h *GRPCHandler) AddClick(ctx context.Context, req *statpb.AddClickRequest) (*statpb.Empty, error) {
	if req.LinkId == 0 {
		return nil, status.Error(codes.InvalidArgument, "link_id is required")
	}
	h.Repo.AddClick(uint32(req.LinkId), uint64(req.UserId))
	return &statpb.Empty{}, nil
}
