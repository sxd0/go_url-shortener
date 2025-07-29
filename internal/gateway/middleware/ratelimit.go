package middleware

import (
	"net"
	"net/http"
	"strings"
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

var rlSkip = []string{"/metrics", "/swagger", "/healthz"}

func (c *RLConfig) Handler(next http.Handler) http.Handler {
	script := redis.NewScript(lua)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := chi.RouteContext(r.Context()).RoutePattern()
		if route == "" {
			parts := strings.SplitN(r.URL.Path, "/", 3)
			if len(parts) > 1 && parts[1] != "" {
				route = "/" + parts[1]
			} else {
				route = "/"
			}
		}

		for _, skip := range rlSkip {
			if strings.HasPrefix(route, skip) {
				next.ServeHTTP(w, r)
				return
			}
		}

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		key := "rl:" + route + ":ip:" + ip

		curr, _ := script.Run(r.Context(), c.Redis,
			[]string{key}, int(c.TTL.Seconds())).Int()

		if curr > c.Limit {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
