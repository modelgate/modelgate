package anthropic

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/modelgate/modelgate/internal/runtime/core"
)

// StreamHandler 流式处理
type StreamHandler struct {
	base *Handler
}

func NewStreamHandler(provider string) *StreamHandler {
	return &StreamHandler{
		base: NewHandler(provider),
	}
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
	resp, err := http.DefaultClient.Do(c.HTTPRequest)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("anthropic stream error: %s", b)
	}

	c.HTTPResponse = resp
	return newStreamReceiver(resp.Body), nil
}
