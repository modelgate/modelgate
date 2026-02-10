package core

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type Execution interface {
	Execute(ctx context.Context, c *Context) error
}

type StreamExecution interface {
	Execute(ctx context.Context, c *Context) error
}

type Executor struct {
	handler Handler
	hooks   []Hook
}

func NewExecutor(h Handler, hooks ...Hook) *Executor {
	return &Executor{
		handler: h,
		hooks:   hooks,
	}
}

// Execute 执行
func (e *Executor) Execute(ctx context.Context, c *Context) (err error) {
	// hooks before
	for _, h := range e.hooks {
		log.Debugf("hook %s before request...", h.Name())
		if hErr := h.Before(ctx, c); hErr != nil {
			log.Errorf("hook %s before error: %v", h.Name(), hErr)
		}
	}

	// execute
	err = e.execute(ctx, c)
	if err != nil {
		c.LastErr = err
	}

	// hooks after（反向）
	for i := len(e.hooks) - 1; i >= 0; i-- {
		log.Debugf("hook %s after request...", e.hooks[i].Name())
		if hErr := e.hooks[i].After(ctx, c); hErr != nil {
			log.Errorf("hook %s after error: %v", e.hooks[i].Name(), hErr)
		}
	}
	return
}

func (e *Executor) execute(ctx context.Context, c *Context) (err error) {
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
