package middleware

import (
    "fmt"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    chimw "github.com/go-chi/chi/v5/middleware"
    "github.com/prometheus/client_golang/prometheus"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{Name: "http_requests_total", Help: "Total HTTP requests"},
        []string{"method", "route", "status_class"},
    )
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "HTTP latency"},
        []string{"method", "route"},
    )
)

func init() { prometheus.MustRegister(httpRequestsTotal, httpRequestDuration) }

func Prometheus(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rw := chimw.NewWrapResponseWriter(w, r.ProtoMajor)
        next.ServeHTTP(rw, r)
        statusClass := fmt.Sprintf("%dxx", rw.Status()/100)
        route := chi.RouteContext(r.Context()).RoutePattern()
        if route == "" { route = "unknown" }
        httpRequestsTotal.WithLabelValues(r.Method, route, statusClass).Inc()
        httpRequestDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
    })
}