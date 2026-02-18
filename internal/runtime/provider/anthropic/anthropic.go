package anthropic

import (
	"github.com/samber/do/v2"

	"github.com/modelgate/modelgate/internal/runtime/core"
	"github.com/modelgate/modelgate/internal/runtime/hooks"
	"github.com/modelgate/modelgate/internal/runtime/provider"
)

func Init(i do.Injector) {
	provider.RegisterPlanSet(core.ProviderCodeAnthropic, &provider.ProviderPlanSet{
		Sync: &provider.SyncExecution{
			Handler: NewHandler(core.ProviderCodeAnthropic),
			Hooks: []core.Hook{
				do.MustInvoke[*hooks.RequestHook](i),
				do.MustInvoke[*hooks.OpenAITokenHook](i),
				do.MustInvoke[*hooks.BillingHook](i),
			},
			Retry: 3,
		},
		Stream: &provider.StreamExecution{
			Handler: NewStreamHandler(core.ProviderCodeAnthropic),
			Hooks: []core.Hook{
				do.MustInvoke[*hooks.RequestHook](i),
				do.MustInvoke[*hooks.StreamWriteHook](i),
				do.MustInvoke[*hooks.OpenAITokenHook](i),
				do.MustInvoke[*hooks.BillingHook](i),
			},
		},
	})
}
