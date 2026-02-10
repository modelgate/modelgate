package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"github.com/modelgate/modelgate/internal/config"
)

// IPRateLimiter 基于 IP 的限流器
type IPRateLimiter struct {
	limiters        sync.Map // map[string]*rateLimiterEntry
	rate            rate.Limit
	burst           int
	cleanupInterval time.Duration
	expireAfter     time.Duration
	stopCleanup     chan struct{}
}

// rateLimiterEntry 限流器条目
type rateLimiterEntry struct {
	limiter    *rate.Limiter
	lastAccess time.Time
	mu         sync.RWMutex
}

// NewIPRateLimiter 创建 IP 限流器
func NewIPRateLimiter(r rate.Limit, b int, cleanupInterval, expireAfter time.Duration) *IPRateLimiter {
	rl := &IPRateLimiter{
		rate:            r,
		burst:           b,
		cleanupInterval: cleanupInterval,
		expireAfter:     expireAfter,
		stopCleanup:     make(chan struct{}),
	}

	rl.startCleanup()
	return rl
}

// GetLimiter 获取或创建指定 IP 的限流器
func (rl *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	// 尝试从 map 获取
	if val, ok := rl.limiters.Load(ip); ok {
		entry := val.(*rateLimiterEntry)
		entry.mu.Lock()
		entry.lastAccess = time.Now()
		entry.mu.Unlock()
		return entry.limiter
	}

	// 创建新的限流器
	limiter := rate.NewLimiter(rl.rate, rl.burst)
	entry := &rateLimiterEntry{
		limiter:    limiter,
		lastAccess: time.Now(),
	}

	// 存储（LoadOrStore 防止并发创建）
	actual, _ := rl.limiters.LoadOrStore(ip, entry)
	return actual.(*rateLimiterEntry).limiter
}

// startCleanup 启动后台清理 goroutine
func (rl *IPRateLimiter) startCleanup() {
	go func() {
		ticker := time.NewTicker(rl.cleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				rl.cleanup()
			case <-rl.stopCleanup:
				return
			}
		}
	}()
}

// cleanup 清理过期的限流器
func (rl *IPRateLimiter) cleanup() {
	now := time.Now()
	cleaned := 0
	remaining := 0

	rl.limiters.Range(func(key, value interface{}) bool {
		entry := value.(*rateLimiterEntry)

		entry.mu.RLock()
		expired := now.Sub(entry.lastAccess) > rl.expireAfter
		entry.mu.RUnlock()

		if expired {
			rl.limiters.Delete(key)
			cleaned++
		} else {
			remaining++
		}

		return true
	})

	if cleaned > 0 {
		log.Infof("Rate limiter cleanup: removed %d expired IPs, %d remaining", cleaned, remaining)
	}
}

// Stop 停止清理 goroutine
func (rl *IPRateLimiter) Stop() {
	close(rl.stopCleanup)
}

// getRealClientIP 提取客户端真实 IP
func getRealClientIP(c *gin.Context, trustProxy bool) string {
	if trustProxy {
		// 优先使用 X-Real-IP
		if ip := c.GetHeader("X-Real-IP"); ip != "" && isValidIP(ip) {
			return ip
		}

		// 其次使用 X-Forwarded-For 的第一个 IP
		if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
			ips := strings.Split(xff, ",")
			if ip := strings.TrimSpace(ips[0]); isValidIP(ip) {
				return ip
			}
		}
	}

	// 降级到 RemoteAddr
	ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
	return ip
}

// isValidIP 验证 IP 格式
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// RateLimit 创建限流中间件
func RateLimit(i do.Injector) gin.HandlerFunc {
	cfg := do.MustInvoke[*config.Config](i)

	// 如果未启用，返回空中间件
	if !cfg.RateLimit.Enabled {
		log.Info("Rate limiter is disabled")
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// 创建 IP 限流器
	limiter := NewIPRateLimiter(
		rate.Limit(cfg.RateLimit.RequestsPerSecond),
		cfg.RateLimit.Burst,
		time.Duration(cfg.RateLimit.CleanupInterval)*time.Second,
		time.Duration(cfg.RateLimit.ExpireAfter)*time.Second,
	)

	log.Infof("Rate limiter enabled: %.0f req/s, burst: %d", cfg.RateLimit.RequestsPerSecond, cfg.RateLimit.Burst)

	return func(c *gin.Context) {
		ip := getRealClientIP(c, cfg.RateLimit.TrustProxy)

		// 获取该 IP 的限流器
		ipLimiter := limiter.GetLimiter(ip)

		// 检查是否允许
		if !ipLimiter.Allow() {
			log.Warnf("Rate limit exceeded for IP: %s, path: %s", ip, c.Request.URL.Path)

			c.Header("Retry-After", "60")
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%.0f", cfg.RateLimit.RequestsPerSecond))
			c.Header("X-RateLimit-Remaining", "0")

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"message": "Rate limit exceeded. Please retry after 60 seconds.",
					"type":    "rate_limit_error",
					"code":    "rate_limit_exceeded",
				},
			})
			return
		}

		c.Next()
	}
}
