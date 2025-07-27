package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type limiterEntry struct {
	limiter *rate.Limiter
	last    time.Time
}

var (
	rps        = 3
	burst      = 6
	limiters   = make(map[string]*limiterEntry)
	limitersMu sync.Mutex
)

func getLimiter(ip string) *rate.Limiter {
	now := time.Now()
	limitersMu.Lock()
	defer limitersMu.Unlock()

	if e, ok := limiters[ip]; ok {
		e.last = now
		return e.limiter
	}
	lim := rate.NewLimiter(rate.Limit(rps), burst)
	limiters[ip] = &limiterEntry{limiter: lim, last: now}
	return lim
}

func cleanupOldLimiters(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for range ticker.C {
		limitersMu.Lock()
		now := time.Now()
		for ip, e := range limiters {
			if now.Sub(e.last) > d {
				delete(limiters, ip)
			}
		}
		limitersMu.Unlock()
	}
}

func init() {
	go cleanupOldLimiters(5 * time.Minute)
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if ip == "" {
			ip = r.RemoteAddr
		}
		lim := getLimiter(ip)
		if !lim.Allow() {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
