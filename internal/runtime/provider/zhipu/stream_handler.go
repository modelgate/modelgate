package zhipu

import (
	"context"

	"github.com/modelgate/modelgate/internal/runtime/core"
	"github.com/modelgate/modelgate/internal/runtime/provider/anthropic"
	"github.com/modelgate/modelgate/internal/runtime/provider/openai"
)

// StreamHandler 流式处理
type StreamHandler struct {
	base         *Handler
	streamHandle core.StreamHandler
}

func NewStreamHandler(provider string) *StreamHandler {
	h := &StreamHandler{
		base: NewHandler(provider),
	}
	if provider == core.ProviderCodeZhipuClaude {
		h.streamHandle = anthropic.NewStreamHandler(provider)
	} else {
		h.streamHandle = openai.NewStreamHandler(provider)
	}
	return h
}

// Provider 获取提供者
func (h StreamHandler) Provider() string {
	return h.base.Provider()
}

// BeforeRequest 构建请求参数
func (h StreamHandler) BeforeRequest(ctx context.Context, c *core.Context) error {
	return h.base.BeforeRequest(ctx, c)
}

// DoStream 发送请求，并处理结果
func (h StreamHandler) DoStream(ctx context.Context, c *core.Context) (stream core.Stream, err error) {
	return h.streamHandle.DoStream(ctx, c)
}
