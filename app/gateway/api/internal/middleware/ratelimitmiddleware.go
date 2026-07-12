package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type limiterItem struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

type RateLimitMiddleware struct {
	rate       rate.Limit
	burst      int
	mu         sync.Mutex
	buckets    map[string]*limiterItem
	maxEntries int
}

func NewRateLimitMiddleware(r, burst int) *RateLimitMiddleware {
	if r <= 0 {
		r = 50
	}
	if burst <= 0 {
		burst = 100
	}
	m := &RateLimitMiddleware{
		rate:       rate.Limit(r),
		burst:      burst,
		buckets:    make(map[string]*limiterItem),
		maxEntries: 10000,
	}
	go m.cleanup()
	return m
}

func (m *RateLimitMiddleware) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for ip, item := range m.buckets {
			if now.Sub(item.lastAccess) > 10*time.Minute {
				delete(m.buckets, ip)
			}
		}
		m.mu.Unlock()
	}
}

func (m *RateLimitMiddleware) limiter(ip string) *rate.Limiter {
	m.mu.Lock()
	defer m.mu.Unlock()
	item, ok := m.buckets[ip]
	if !ok {
		if len(m.buckets) >= m.maxEntries {
			var oldestIP string
			var oldestTime time.Time
			for k, v := range m.buckets {
				if oldestIP == "" || v.lastAccess.Before(oldestTime) {
					oldestIP = k
					oldestTime = v.lastAccess
				}
			}
			if oldestIP != "" {
				delete(m.buckets, oldestIP)
			}
		}
		item = &limiterItem{
			limiter:    rate.NewLimiter(m.rate, m.burst),
			lastAccess: time.Now(),
		}
		m.buckets[ip] = item
	}
	item.lastAccess = time.Now()
	return item.limiter
}

func (m *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		if !m.limiter(ip).Allow() {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"code":429,"msg":"too many requests"}`))
			return
		}
		next(w, r)
	}
}
