package openai

import (
	"github.com/samber/do/v2"

	"github.com/modelgate/modelgate/internal/runtime/core"
	"github.com/modelgate/modelgate/internal/runtime/hooks"
	"github.com/modelgate/modelgate/internal/runtime/provider"
)

func Init(i do.Injector) {
	baseHooks := []core.Hook{
		do.MustInvoke[*hooks.RequestHook](i),
		do.MustInvoke[*hooks.OpenAITokenHook](i),
		do.MustInvoke[*hooks.BillingHook](i),
	}

	streamHooks := []core.Hook{
		do.MustInvoke[*hooks.RequestHook](i),
		do.MustInvoke[*hooks.StreamWriteHook](i),
		do.MustInvoke[*hooks.OpenAITokenHook](i),
		do.MustInvoke[*hooks.BillingHook](i),
	}

	// OpenAI
	provider.RegisterPlanSet((core.ProviderCodeOpenAI), &provider.ProviderPlanSet{
		Sync: &provider.SyncExecution{
			Handler: NewHandler(core.ProviderCodeOpenAI),
			Hooks:   baseHooks,
			Retry:   3,
		},
		Stream: &provider.StreamExecution{
			Handler: NewStreamHandler(core.ProviderCodeOpenAI),
			Hooks:   streamHooks,
		},
	})

	// DeepSeek
	provider.RegisterPlanSet(core.ProviderCodeDeepSeek, &provider.ProviderPlanSet{
		Sync: &provider.SyncExecution{
			Handler: NewHandler(core.ProviderCodeDeepSeek),
			Hooks:   baseHooks,
			Retry:   3,
		},
		Stream: &provider.StreamExecution{
			Handler: NewStreamHandler(core.ProviderCodeDeepSeek),
			Hooks:   streamHooks,
		},
	})
}
