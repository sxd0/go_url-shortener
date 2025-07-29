package middleware

import (
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type RLConfig struct {
	Limit int
	TTL   time.Duration
	Redis *redis.Client
}

const lua = `
local current = redis.call("INCR", KEYS[1])
if current == 1 then
	redis.call("EXPIRE", KEYS[1], ARGV[1])
end
return current
`

func (c *RLConfig) Handler(next http.Handler) http.Handler {
	script := redis.NewScript(lua)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := chi.RouteContext(r.Context()).RoutePattern()
		if route == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}
		if route == "" {
			route = "unknown"
		}

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		key := "rl:" + route + ":ip:" + ip

		curr, _ := script.Run(r.Context(), c.Redis, []string{key}, int(c.TTL.Seconds())).Int()
		if curr > c.Limit {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
