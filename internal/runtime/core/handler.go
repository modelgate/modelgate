package core

import (
	"context"
)

// Handler 处理器接口（统一非流式和流式处理）
type Handler interface {
	Provider() string
	BeforeRequest(ctx context.Context, c *Context) error
	DoRequest(ctx context.Context, c *Context) error
	AfterResponse(ctx context.Context, c *Context) error
	DoStream(ctx context.Context, c *Context) (Stream, error)
}
