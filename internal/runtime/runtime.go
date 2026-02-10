package runtime

import (
	"context"

	"github.com/samber/do/v2"
	log "github.com/sirupsen/logrus"

	"github.com/modelgate/modelgate/internal/runtime/core"
	"github.com/modelgate/modelgate/internal/runtime/hooks"
	"github.com/modelgate/modelgate/internal/runtime/provider"
	"github.com/modelgate/modelgate/internal/runtime/provider/anthropic"
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
	zhipu.Init(i)

	do.Provide(i, provider.NewProviderPlanner)
	do.Provide(i, New)
}

// New 实例化
func New(i do.Injector) (*Runtime, error) {
	planner := do.MustInvoke[provider.Planner](i)
	return &Runtime{
		planner: planner,
	}, nil
}

// Runtime 运行时
type Runtime struct {
	planner provider.Planner
}

// Plan 计划
func (r *Runtime) Plan(c *core.Context) ([]*provider.Plan, error) {
	log.Infof("is stream: %v", c.IsStream)
	return r.planner.Plan(c)
}

// Run 运行
func (r *Runtime) Run(ctx context.Context, c *core.Context) error {
	task := &Task{
		Ctx:     ctx,
		Context: c,
	}
	return task.Run(r)
}
