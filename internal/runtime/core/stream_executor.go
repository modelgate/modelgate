package core

import (
	"context"
	"io"

	log "github.com/sirupsen/logrus"
)

type streamExecutor struct {
	handler Handler
	hooks   []Hook
}

// NewStreamExecutor 创建流式执行器
func NewStreamExecutor(h Handler, hooks ...Hook) Executor {
	return &streamExecutor{
		handler: h,
		hooks:   hooks,
	}
}

// Execute 执行流式处理
func (e *streamExecutor) Execute(ctx context.Context, c *Context) error {
	// hooks before
	for _, h := range e.hooks {
		log.Debugf("hook: %s, before...", h.Name())
		if hErr := h.Before(ctx, c); hErr != nil {
			log.Errorf("hook %s before error: %v", h.Name(), hErr)
		}
	}

	// before request
	log.Debugf("provider %s, model: %s, before request...", e.handler.Provider(), c.CurrentModel.ModelCode)
	err := e.handler.BeforeRequest(ctx, c)
	if err != nil {
		log.Errorf("provider %s before request error: %v", e.handler.Provider(), err)
		return err
	}

	// do request
	log.Debugf("provider %s, model: %s, do request...", e.handler.Provider(), c.CurrentModel.ModelCode)
	stream, err := e.handler.DoStream(ctx, c)
	if err != nil {
		log.Errorf("provider %s do stream error: %v", e.handler.Provider(), err)
		return err
	}
	defer stream.Close()

	for {
		chunk, sErr := stream.Recv()
		if sErr != nil && sErr != io.EOF {
			e.callOnError(ctx, c, sErr)
			return sErr
		}

		for _, h := range e.hooks {
			if err := h.OnChunk(ctx, c, chunk); err != nil {
				log.Errorf("hook %s on chunk error: %v", h.Name(), err)
				return err
			}
		}

		if sErr == io.EOF {
			break
		}
	}

	// hooks after（反向）
	for i := len(e.hooks) - 1; i >= 0; i-- {
		log.Debugf("after hook %s ...", e.hooks[i].Name())
		if hErr := e.hooks[i].After(ctx, c); hErr != nil {
			log.Errorf("hook %s after error: %v", e.hooks[i].Name(), hErr)
		}
	}
	return nil
}

func (e *streamExecutor) callOnError(ctx context.Context, c *Context, err error) {
	for _, h := range e.hooks {
		h.OnError(ctx, c, err)
	}
}
