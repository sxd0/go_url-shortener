package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	goredis "github.com/redis/go-redis/v9"
)

type RLConfig struct {
	Limit          int
	TTL            time.Duration
	Redis          *goredis.Client
	TrustedProxies []string
	KeyMode        string // "global" | "route" | "route+hash"
}

var rlSkip = []string{"/metrics", "/swagger", "/healthz"}

func (c *RLConfig) Handler(next http.Handler) http.Handler {
	if c == nil || c.Redis == nil || c.Limit <= 0 || c.TTL <= 0 {
		return next
	}

	trustedNets := parseCIDRs(c.TrustedProxies)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := routeKey(r, c.KeyMode)

		for _, skip := range rlSkip {
			if strings.HasPrefix(route, skip) {
				next.ServeHTTP(w, r)
				return
			}
		}

		ip := clientIP(r, trustedNets)
		key := "rl:v1:" + route + ":ip:" + ip

		nowMs := time.Now().UnixMilli()
		oldestMs := nowMs - c.TTL.Milliseconds()

		ctx, cancel := context.WithTimeout(r.Context(), 100*time.Millisecond)
		defer cancel()

		member := fmt.Sprintf("%d-%d", nowMs, rand.Int63())

		pipe := c.Redis.TxPipeline()
		pipe.ZAdd(ctx, key, goredis.Z{Score: float64(nowMs), Member: member})
		pipe.ZRemRangeByScore(ctx, key, "-inf", fmtScore(oldestMs))
		count := pipe.ZCard(ctx, key)
		pipe.Expire(ctx, key, c.TTL)

		_, err := pipe.Exec(ctx)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		n, _ := count.Result()
		if int(n) > c.Limit {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func fmtScore(ms int64) string {
	return fmt.Sprintf("%f", float64(ms))
}

func routeKey(r *http.Request, mode string) string {
	switch mode {
	case "global":
		return "global"
	case "route":
		if rc := chi.RouteContext(r.Context()); rc != nil && rc.RoutePattern() != "" {
			return rc.RoutePattern()
		}
		if r.URL.Path == "/" {
			return "/"
		}
		parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 2)
		return "/" + parts[0]
	case "route+hash":
		if rc := chi.RouteContext(r.Context()); rc != nil && rc.RoutePattern() != "" {
			if strings.HasPrefix(rc.RoutePattern(), "/r/") {
				h := chi.URLParam(r, "hash")
				return "/r/" + h
			}
			return rc.RoutePattern()
		}
		return r.URL.Path
	default:
		if rc := chi.RouteContext(r.Context()); rc != nil && rc.RoutePattern() != "" {
			return rc.RoutePattern()
		}
		return r.URL.Path
	}
}

func clientIP(r *http.Request, trusted []*net.IPNet) string {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	remote := net.ParseIP(host)

	if remote != nil && ipInNets(remote, trusted) {
		if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
			parts := strings.Split(xf, ",")
			first := strings.TrimSpace(parts[0])
			if ip := net.ParseIP(first); ip != nil {
				return first
			}
		}
		if xr := strings.TrimSpace(r.Header.Get("X-Real-IP")); xr != "" {
			if ip := net.ParseIP(xr); ip != nil {
				return xr
			}
		}
	}
	if host != "" {
		return host
	}
	return "unknown"
}

func parseCIDRs(list []string) []*net.IPNet {
	if len(list) == 0 {
		list = []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.1/32", "::1/128"}
	}
	var nets []*net.IPNet
	for _, s := range list {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ipn, err := net.ParseCIDR(s); err == nil {
			nets = append(nets, ipn)
		}
	}
	return nets
}

func ipInNets(ip net.IP, nets []*net.IPNet) bool {
	for _, n := range nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}
