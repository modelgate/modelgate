package core

import "context"

// Hook 非流式处理
type Hook interface {
	Name() string
	Before(ctx context.Context, c *Context) error
	After(ctx context.Context, c *Context) error
}

// StreamHook 流式处理
type StreamHook interface {
	Name() string
	BeforeStream(ctx context.Context, c *Context) error
	OnChunk(ctx context.Context, c *Context, chunk *StreamChunk) error
	AfterStream(ctx context.Context, c *Context) error
	OnError(ctx context.Context, c *Context, err error)
}
