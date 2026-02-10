package runtime

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/modelgate/modelgate/internal/runtime/core"
)

// Task 任务
type Task struct {
	Ctx     context.Context
	Context *core.Context // modelgate 请求上下文

	LastErr error
}

// Run 运行
func (t *Task) Run(r *Runtime) error {
	plans, err := r.Plan(t.Context)
	if err != nil {
		return err
	}

	log.Infof("plans: %d", len(plans))
	for _, plan := range plans {
		if err = plan.Exec.Execute(t.Ctx, t.Context); err == nil {
			return nil
		}
		t.LastErr = err
		log.Errorf("plan execute error: %v", err)
	}
	if t.LastErr != nil {
		return fmt.Errorf("task failed: %v", t.LastErr)
	}
	return errors.New("task failed")
}
