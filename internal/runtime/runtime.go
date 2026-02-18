package runtime

import (
	"context"

	"github.com/samber/do/v2"

	"github.com/modelgate/modelgate/internal/runtime/core"
	"github.com/modelgate/modelgate/internal/runtime/hooks"
	"github.com/modelgate/modelgate/internal/runtime/provider/anthropic"
	"github.com/modelgate/modelgate/internal/runtime/provider/minimax"
	"github.com/modelgate/modelgate/internal/runtime/provider/openai"
	"github.com/modelgate/modelgate/internal/runtime/provider/zhipu"
)

// Init 初始化
func Init(i do.Injector) {
	// Hooks
	do.Provide(i, hooks.NewRequestHook)
	do.Provide(i, hooks.NewStreamHook)
	do.Provide(i, hooks.NewOpenAITokenHook)
	do.Provide(i, hooks.NewBillingHook)

	// Provider
	anthropic.Init(i)
	openai.Init(i)
	minimax.Init(i)
	zhipu.Init(i)
}

// Run 执行
func Run(ctx context.Context, c *core.Context) (err error) {
	exector, err := core.ExecutorRegistry.Get(c.CurrentModel.ProviderCode, core.Options{
		IsStream:    c.IsStream,
		IsAnthropic: c.IsAnthropic,
		Retry:       3,
	})
	if err != nil {
		return
	}

	err = exector.Execute(ctx, c)
	return
}
