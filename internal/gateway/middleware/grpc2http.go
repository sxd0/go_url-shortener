package middleware

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WriteGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	msg := st.Message()
	switch st.Code() {
	case codes.InvalidArgument:
		http.Error(w, msg, http.StatusBadRequest)
	case codes.Unauthenticated:
		http.Error(w, msg, http.StatusUnauthorized)
	case codes.PermissionDenied:
		http.Error(w, msg, http.StatusForbidden)
	case codes.NotFound:
		http.Error(w, msg, http.StatusNotFound)
	case codes.AlreadyExists:
		http.Error(w, msg, http.StatusConflict)
	case codes.ResourceExhausted:
		http.Error(w, msg, http.StatusTooManyRequests)
	case codes.Unavailable:
		http.Error(w, msg, http.StatusServiceUnavailable)
	case codes.DeadlineExceeded:
		http.Error(w, msg, http.StatusGatewayTimeout)
	default:
		http.Error(w, msg, http.StatusInternalServerError)
	}
}
