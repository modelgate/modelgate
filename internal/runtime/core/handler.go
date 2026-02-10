package core

import (
	"context"
)

// Handler 非流式处理
type Handler interface {
	Provider() string
	BeforeRequest(ctx context.Context, c *Context) error
	DoRequest(ctx context.Context, c *Context) error
	AfterResponse(ctx context.Context, c *Context) error
}

// StreamHandler 流式处理
type StreamHandler interface {
	Provider() string
	BeforeRequest(ctx context.Context, c *Context) error
	DoStream(ctx context.Context, c *Context) (Stream, error)
}
