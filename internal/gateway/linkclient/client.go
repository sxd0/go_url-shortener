package linkclient

import (
	"context"
	"time"

	"github.com/sxd0/go_url-shortener/proto/gen/go/linkpb"
	"google.golang.org/grpc"
)

type LinkClient struct {
	client linkpb.LinkServiceClient
}

func NewLinkClient(conn *grpc.ClientConn) *LinkClient {
	return &LinkClient{
		client: linkpb.NewLinkServiceClient(conn),
	}
}

func (l *LinkClient) CreateLink(ctx context.Context, in *linkpb.CreateLinkRequest, opts ...grpc.CallOption) (*linkpb.LinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return l.client.CreateLink(ctx, in)
}

func (l *LinkClient) GetAllLinks(ctx context.Context, in *linkpb.GetAllLinksRequest, opts ...grpc.CallOption) (*linkpb.GetAllLinksResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return l.client.GetAllLinks(ctx, in)
}

func (l *LinkClient) GetLinkByHash(ctx context.Context, in *linkpb.GetLinkByHashRequest, opts ...grpc.CallOption) (*linkpb.LinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return l.client.GetLinkByHash(ctx, in)
}

func (l *LinkClient) UpdateLink(ctx context.Context, in *linkpb.UpdateLinkRequest, opts ...grpc.CallOption) (*linkpb.LinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return l.client.UpdateLink(ctx, in)
}

func (l *LinkClient) DeleteLink(ctx context.Context, in *linkpb.DeleteLinkRequest, opts ...grpc.CallOption) (*linkpb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return l.client.DeleteLink(ctx, in)
}
