package middleware

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

func WithAuthMD(ctx context.Context, r *http.Request) context.Context {
	if v := r.Header.Get("Authorization"); v != "" {
		return metadata.AppendToOutgoingContext(ctx, "authorization", v)
	}
	return ctx
}
