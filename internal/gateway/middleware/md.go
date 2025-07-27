package middleware

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

func AttachCommonMD(ctx context.Context, r *http.Request) context.Context {
	// Authorization
	if v := r.Header.Get("Authorization"); v != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", v)
	}
	// X-Request-ID
	if rid := GetRequestIDFromContext(ctx); rid != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", rid)
	}
	return ctx
}
