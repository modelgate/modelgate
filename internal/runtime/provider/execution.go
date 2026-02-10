package provider

import (
	"context"
	"errors"

	"github.com/modelgate/modelgate/internal/runtime/core"
)

// Execution 执行器
type Execution interface {
	Execute(ctx context.Context, c *core.Context) error
}

// SyncExecution 同步执行
type SyncExecution struct {
	Handler core.Handler
	Hooks   []core.Hook
	Retry   int
}

// Execute 执行
func (e *SyncExecution) Execute(ctx context.Context, c *core.Context) error {
	executor := core.NewExecutor(e.Handler, e.Hooks...)
	for i := 0; i <= e.Retry; i++ {
		c.AttemptNo = i + 1
		c.LastErr = nil

		if err := executor.Execute(ctx, c); err == nil {
			return nil
		}
	}
	return errors.New("execute failed")
}

// StreamExecution 流式执行
type StreamExecution struct {
	Handler core.StreamHandler
	Hooks   []core.StreamHook
}

// Execute 执行
func (e *StreamExecution) Execute(ctx context.Context, c *core.Context) error {
	c.AttemptNo = 1
	stream := core.NewStreamExecutor(e.Handler, e.Hooks...)
	return stream.Execute(ctx, c)
}
