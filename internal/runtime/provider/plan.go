package provider

import (
	"errors"
	"sync"

	"github.com/modelgate/modelgate/internal/runtime/core"
	"github.com/samber/do/v2"
)

var plansMap sync.Map

func RegisterPlanSet(provider string, planSet *ProviderPlanSet) {
	plansMap.Store(provider, planSet)
}

type ProviderPlanSet struct {
	Sync   Execution
	Stream Execution
}

// Planner 计划器
type Planner interface {
	Plan(ctx *core.Context) ([]*Plan, error)
}

// Plan 计划
type Plan struct {
	Exec Execution
}

type ProviderPlanner struct {
}

func NewProviderPlanner(i do.Injector) (Planner, error) {
	return &ProviderPlanner{}, nil
}

// ProviderPlanner 根据提供者
func (p *ProviderPlanner) Plan(c *core.Context) (plans []*Plan, err error) {
	planSet, ok := plansMap.Load(c.CurrentModel.ProviderCode)
	if !ok {
		err = errors.New("provider not found")
		return
	}
	if c.IsStream {
		plans = append(plans, &Plan{Exec: planSet.(*ProviderPlanSet).Stream})
	} else {
		plans = append(plans, &Plan{Exec: planSet.(*ProviderPlanSet).Sync})
	}
	return
}
