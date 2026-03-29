package gemini

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

	// Gemini 支持 OpenAI 协议
	core.ExecutorRegistry.Register(core.ProviderCodeGemini, func(opts core.Options) (core.Executor, error) {
		if opts.IsStream {
			return core.NewStreamExecutor(openaiHandler, reqHook, streamWriteHook, tokenHook, billingHook), nil
		}
		base := core.NewExecutor(openaiHandler, reqHook, tokenHook, billingHook)
		return core.NewRetryExecutor(base, opts.Retry), nil
	})
}