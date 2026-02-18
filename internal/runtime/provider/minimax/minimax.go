package minimax

import (
	"github.com/samber/do/v2"

	"github.com/modelgate/modelgate/internal/runtime/core"
	"github.com/modelgate/modelgate/internal/runtime/hooks"
)

func Init(i do.Injector) {
	reqHook := do.MustInvoke[*hooks.RequestHook](i)
	tokenHook := do.MustInvoke[*hooks.OpenAITokenHook](i)
	billingHook := do.MustInvoke[*hooks.BillingHook](i)
	streamWriteHook := do.MustInvoke[*hooks.StreamWriteHook](i)

	openaiHandler := NewOpenAIHandler()
	anthropicHandler := NewAnthropicHandler()

	// MiniMax 同时支持 OpenAI 和 Anthropic 协议，根据 opts.UrlPath 选择对应 handler
	core.ExecutorRegistry.Register(core.ProviderCodeMinimax, func(opts core.Options) (core.Executor, error) {
		var handler core.Handler
		if opts.IsAnthropic {
			handler = anthropicHandler
		} else {
			handler = openaiHandler
		}

		if opts.IsStream {
			return core.NewStreamExecutor(handler, reqHook, streamWriteHook, tokenHook, billingHook), nil
		}
		base := core.NewExecutor(handler, reqHook, tokenHook, billingHook)
		return core.NewRetryExecutor(base, opts.Retry), nil
	})
}
