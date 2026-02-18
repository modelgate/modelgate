package core

import (
	"net/http"
	"time"
)

// HttpClient 全局 HTTP 客户端
var HttpClient = &http.Client{
	Timeout: 300 * time.Second,
}

// DefaultHttpClient 返回默认客户端（用于需要 *http.Client 的场景）
func DefaultHttpClient() *http.Client {
	return HttpClient
}
