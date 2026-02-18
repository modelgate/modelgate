package hooks

import (
	"context"

	"github.com/samber/do/v2"

	"github.com/modelgate/modelgate/internal/runtime/core"
)

type StreamWriteHook struct {
}

var _ core.Hook = (*StreamWriteHook)(nil)

func NewStreamHook(i do.Injector) (*StreamWriteHook, error) {
	return &StreamWriteHook{}, nil
}

func (h *StreamWriteHook) Name() string {
	return "stream_write"
}

func (h *StreamWriteHook) Before(ctx context.Context, c *core.Context) error {
	if c.StreamWriter != nil {
		_ = c.StreamWriter.Open()
	}
	return nil
}

func (h *StreamWriteHook) After(ctx context.Context, c *core.Context) error {
	if c.StreamWriter != nil {
		_ = c.StreamWriter.Close()
	}
	return nil
}

func (h *StreamWriteHook) OnChunk(ctx context.Context, c *core.Context, chunk *core.StreamChunk) error {
	if c.StreamWriter != nil {
		_ = c.StreamWriter.Write(chunk)
	}
	return nil
}

func (h *StreamWriteHook) OnError(ctx context.Context, c *core.Context, err error) {
	if c.StreamWriter != nil {
		_ = c.StreamWriter.Close()
	}
}
