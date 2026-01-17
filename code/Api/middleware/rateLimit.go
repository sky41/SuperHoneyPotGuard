package middleware

import (
	"net/http"
	"superhoneypotguard/config"
	"superhoneypotguard/utils"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*visitor
	mu  sync.RWMutex
	r   rateConfig
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type rateConfig struct {
	limit rate.Limit
	burst int
}

func NewIPRateLimiter(r rateConfig) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*visitor),
		r:   r,
	}
}

func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	v, exists := i.ips[ip]
	if !exists {
		limiter := rate.NewLimiter(i.r.limit, i.r.burst)
		i.ips[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (i *IPRateLimiter) CleanupStaleVisitors() {
	i.mu.Lock()
	defer i.mu.Unlock()

	for ip, v := range i.ips {
		if time.Since(v.lastSeen) > 3*time.Minute {
			delete(i.ips, ip)
		}
	}
}

var globalLimiter *IPRateLimiter

func InitRateLimiter() {
	cfg := config.AppConfig
	window := time.Duration(cfg.RateLimitWindow) * time.Millisecond
	maxRequests := cfg.RateLimitMax

	limiter := rate.Every(time.Duration(window.Milliseconds()/int64(maxRequests)) * time.Millisecond)
	globalLimiter = NewIPRateLimiter(rateConfig{
		limit: limiter,
		burst: maxRequests,
	})

	go func() {
		for range time.Tick(time.Minute) {
			globalLimiter.CleanupStaleVisitors()
		}
	}()
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := utils.GetClientIP(c)
		limiter := globalLimiter.getLimiter(ip)

		if !limiter.Allow() {
			utils.ErrorResponse(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := utils.GetClientIP(c)
		limiter := rate.NewLimiter(rate.Every(3*time.Minute), 5)

		if !limiter.Allow() {
			utils.ErrorResponse(c, http.StatusTooManyRequests, "登录请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
