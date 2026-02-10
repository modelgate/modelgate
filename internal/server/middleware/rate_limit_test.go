package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func TestGetRealClientIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		remoteAddr string
		headers    map[string]string
		trustProxy bool
		want       string
	}{
		{
			name:       "X-Real-IP with trust proxy",
			remoteAddr: "127.0.0.1:1234",
			headers: map[string]string{
				"X-Real-IP": "192.168.1.100",
			},
			trustProxy: true,
			want:       "192.168.1.100",
		},
		{
			name:       "X-Forwarded-For with trust proxy",
			remoteAddr: "127.0.0.1:1234",
			headers: map[string]string{
				"X-Forwarded-For": "10.0.0.1, 10.0.0.2, 10.0.0.3",
			},
			trustProxy: true,
			want:       "10.0.0.1",
		},
		{
			name:       "X-Real-IP priority over X-Forwarded-For",
			remoteAddr: "127.0.0.1:1234",
			headers: map[string]string{
				"X-Real-IP":       "192.168.1.100",
				"X-Forwarded-For": "10.0.0.1",
			},
			trustProxy: true,
			want:       "192.168.1.100",
		},
		{
			name:       "RemoteAddr fallback",
			remoteAddr: "203.0.113.1:5678",
			headers:    map[string]string{},
			trustProxy: true,
			want:       "203.0.113.1",
		},
		{
			name:       "Ignore proxy headers when not trusted",
			remoteAddr: "203.0.113.1:5678",
			headers: map[string]string{
				"X-Real-IP":       "192.168.1.100",
				"X-Forwarded-For": "10.0.0.1",
			},
			trustProxy: false,
			want:       "203.0.113.1",
		},
		{
			name:       "Invalid IP in X-Real-IP",
			remoteAddr: "203.0.113.1:5678",
			headers: map[string]string{
				"X-Real-IP": "not-an-ip",
			},
			trustProxy: true,
			want:       "203.0.113.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)
			c.Request.RemoteAddr = tt.remoteAddr

			for k, v := range tt.headers {
				c.Request.Header.Set(k, v)
			}

			got := getRealClientIP(c, tt.trustProxy)
			if got != tt.want {
				t.Errorf("getRealClientIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidIP(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{"Valid IPv4", "192.168.1.1", true},
		{"Valid IPv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", true},
		{"Invalid IP", "not-an-ip", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidIP(tt.ip); got != tt.want {
				t.Errorf("isValidIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPRateLimiter_GetLimiter(t *testing.T) {
	rl := NewIPRateLimiter(rate.Limit(10), 10, 1*time.Second, 1*time.Second)
	defer rl.Stop()

	// 测试获取相同 IP 返回相同限流器
	ip := "192.168.1.1"
	limiter1 := rl.GetLimiter(ip)
	limiter2 := rl.GetLimiter(ip)

	if limiter1 != limiter2 {
		t.Error("Expected same limiter for same IP")
	}

	// 测试不同 IP 返回不同限流器
	ip2 := "192.168.1.2"
	limiter3 := rl.GetLimiter(ip2)

	if limiter1 == limiter3 {
		t.Error("Expected different limiters for different IPs")
	}
}

func TestIPRateLimiter_Allow(t *testing.T) {
	// 创建限流器：2 req/s, burst=2
	rl := NewIPRateLimiter(rate.Limit(2), 2, 1*time.Second, 1*time.Second)
	defer rl.Stop()

	ip := "192.168.1.1"
	limiter := rl.GetLimiter(ip)

	// 前 2 个请求应该通过（burst=2）
	if !limiter.Allow() {
		t.Error("First request should be allowed")
	}
	if !limiter.Allow() {
		t.Error("Second request should be allowed")
	}

	// 第 3 个请求应该被拒绝
	if limiter.Allow() {
		t.Error("Third request should be denied")
	}

	// 等待 0.5 秒后，应该有 1 个令牌可用（2 req/s = 0.5s/req）
	time.Sleep(550 * time.Millisecond)
	if !limiter.Allow() {
		t.Error("Request after waiting should be allowed")
	}
}

func TestIPRateLimiter_IPIsolation(t *testing.T) {
	rl := NewIPRateLimiter(rate.Limit(1), 1, 1*time.Second, 1*time.Second)
	defer rl.Stop()

	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	limiter1 := rl.GetLimiter(ip1)
	limiter2 := rl.GetLimiter(ip2)

	// IP1 消耗令牌
	if !limiter1.Allow() {
		t.Error("IP1 first request should be allowed")
	}
	if limiter1.Allow() {
		t.Error("IP1 second request should be denied")
	}

	// IP2 应该有独立的配额
	if !limiter2.Allow() {
		t.Error("IP2 first request should be allowed")
	}
	if limiter2.Allow() {
		t.Error("IP2 second request should be denied")
	}
}

func TestIPRateLimiter_Cleanup(t *testing.T) {
	// 使用较短的过期时间进行测试
	rl := NewIPRateLimiter(rate.Limit(10), 10, 100*time.Millisecond, 200*time.Millisecond)
	defer rl.Stop()

	// 添加多个 IP
	ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}
	for _, ip := range ips {
		rl.GetLimiter(ip)
	}

	// 验证 IP 已存储
	count := 0
	rl.limiters.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	if count != 3 {
		t.Errorf("Expected 3 IPs, got %d", count)
	}

	// 等待清理触发
	time.Sleep(350 * time.Millisecond)

	// 验证过期 IP 已被清理
	count = 0
	rl.limiters.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	if count != 0 {
		t.Errorf("Expected 0 IPs after cleanup, got %d", count)
	}
}

func TestIPRateLimiter_Concurrent(t *testing.T) {
	rl := NewIPRateLimiter(rate.Limit(100), 100, 1*time.Second, 1*time.Second)
	defer rl.Stop()

	ip := "192.168.1.1"

	// 并发访问同一个 IP
	done := make(chan bool)
	for range 10 {
		go func() {
			for j := 0; j < 10; j++ {
				rl.GetLimiter(ip).Allow()
			}
			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for range 10 {
		<-done
	}

	// 验证只有一个限流器被创建
	count := 0
	rl.limiters.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	if count != 1 {
		t.Errorf("Expected 1 limiter, got %d", count)
	}
}

func BenchmarkIPRateLimiter_GetLimiter(b *testing.B) {
	rl := NewIPRateLimiter(rate.Limit(1000), 1000, 10*time.Second, 10*time.Second)
	defer rl.Stop()

	ip := "192.168.1.1"

	for b.Loop() {
		rl.GetLimiter(ip)
	}
}

func BenchmarkIPRateLimiter_Allow(b *testing.B) {
	rl := NewIPRateLimiter(rate.Limit(10000), 10000, 10*time.Second, 10*time.Second)
	defer rl.Stop()

	ip := "192.168.1.1"
	limiter := rl.GetLimiter(ip)

	b.ResetTimer()
	for b.Loop() {
		limiter.Allow()
	}
}
