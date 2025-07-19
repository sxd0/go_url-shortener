package middleware

import (
	"go/test-http/pkg/logger"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &WrapperWriter{ResponseWriter: w, StatusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		logger.Log.Info("request",
			zap.Int("status", ww.StatusCode),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", time.Since(start)),
		)
	})
}
