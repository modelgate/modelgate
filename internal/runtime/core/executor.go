package core

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type Executor interface {
	Execute(ctx context.Context, c *Context) (err error)
}

type executor struct {
	*baseExecutor
	handler Handler
}

func NewExecutor(h Handler, hooks ...Hook) Executor {
	return &executor{
		baseExecutor: newBaseExecutor(hooks),
		handler:      h,
	}
}

// Execute 执行
func (e *executor) Execute(ctx context.Context, c *Context) (err error) {
	// hooks before
	e.callHooksBefore(ctx, c)

	// execute
	err = e.execute(ctx, c)
	if err != nil {
		c.LastErr = err
	}

	// hooks after（反向）
	e.callHooksAfter(ctx, c)
	return
}

func (e *executor) execute(ctx context.Context, c *Context) (err error) {
	// provider before
	log.Debugf("provider %s, model: %s before request...", e.handler.Provider(), c.CurrentModel.ModelCode)
	if err = e.handler.BeforeRequest(ctx, c); err != nil {
		log.Errorf("provider %s before error: %v", e.handler.Provider(), err)
		return
	}

	// do request
	log.Debugf("provider %s, model: %s do request...", e.handler.Provider(), c.CurrentModel.ModelCode)
	if err = e.handler.DoRequest(ctx, c); err != nil {
		log.Errorf("provider %s do request error: %v", e.handler.Provider(), err)
		return
	}

	// provider after
	log.Debugf("provider %s, model: %s after request...", e.handler.Provider(), c.CurrentModel.ModelCode)
	if err = e.handler.AfterResponse(ctx, c); err != nil {
		log.Errorf("provider %s after error: %v", e.handler.Provider(), err)
		return
	}
	return
}

type baseExecutor struct {
	hooks []Hook
}

func (e *baseExecutor) callHooksBefore(ctx context.Context, c *Context) {
	for _, h := range e.hooks {
		log.Debugf("hook %s before request...", h.Name())
		if hErr := h.Before(ctx, c); hErr != nil {
			log.Errorf("hook %s before error: %v", h.Name(), hErr)
		}
	}
}

func (e *baseExecutor) callHooksAfter(ctx context.Context, c *Context) {
	for i := len(e.hooks) - 1; i >= 0; i-- {
		log.Debugf("hook %s after request...", e.hooks[i].Name())
		if hErr := e.hooks[i].After(ctx, c); hErr != nil {
			log.Errorf("hook %s after error: %v", e.hooks[i].Name(), hErr)
		}
	}
}

func newBaseExecutor(hooks []Hook) *baseExecutor {
	return &baseExecutor{
		hooks: hooks,
	}
}

type retryExecutor struct {
	base  Executor
	retry int
}

// NewRetryExecutor 创建重试执行器，retry <= 0 时返回原执行器
func NewRetryExecutor(base Executor, retry int) Executor {
	if retry <= 0 {
		return base
	}
	return &retryExecutor{
		base:  base,
		retry: retry,
	}
}

// Execute 执行并重试
func (r *retryExecutor) Execute(ctx context.Context, c *Context) (err error) {
	for i := 0; i <= r.retry; i++ {
		c.AttemptNo = i + 1
		c.LastErr = nil
		err = r.base.Execute(ctx, c)
		if err == nil {
			return nil
		}
		log.Warnf("executor attempt %d failed: %v", c.AttemptNo, err)
	}
	return err
}
