package openai

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

	{
		handler := NewHandler(core.ProviderCodeOpenAI)

		core.ExecutorRegistry.Register(core.ProviderCodeOpenAI, func(opts core.Options) (core.Executor, error) {
			if opts.IsStream {
				return core.NewStreamExecutor(handler, reqHook, streamWriteHook, tokenHook, billingHook), nil
			} else {
				base := core.NewExecutor(handler, reqHook, tokenHook, billingHook)
				return core.NewRetryExecutor(base, opts.Retry), nil
			}
		})
	}
	{
		handler := NewHandler(core.ProviderCodeDeepSeek)

		core.ExecutorRegistry.Register(core.ProviderCodeDeepSeek, func(opts core.Options) (core.Executor, error) {
			if opts.IsStream {
				return core.NewStreamExecutor(handler, reqHook, streamWriteHook, tokenHook, billingHook), nil
			} else {
				base := core.NewExecutor(handler, reqHook, tokenHook, billingHook)
				return core.NewRetryExecutor(base, opts.Retry), nil
			}
		})
	}
}
