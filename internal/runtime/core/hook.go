package core

import "context"

// Hook Hook 接口（统一非流式和流式处理）
type Hook interface {
	Name() string
	Before(ctx context.Context, c *Context) error
	After(ctx context.Context, c *Context) error
	OnChunk(ctx context.Context, c *Context, chunk *StreamChunk) error
	OnError(ctx context.Context, c *Context, err error)
}
