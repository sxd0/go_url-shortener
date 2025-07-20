package middleware

import (
	"net/http"
	"time"

	"github.com/sxd0/go_url-shortener/pkg/logger"

	"go.uber.org/zap"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := GetRequestID(r.Context())
		ww := &WrapperWriter{ResponseWriter: w, StatusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		logger.Log.Info("request",
			zap.String("request_id", requestID),
			zap.Int("status", ww.StatusCode),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", time.Since(start)),
		)
	})
}
