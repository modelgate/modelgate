package core

import (
	"fmt"
	"sync"
)

type NewExecutorFunc func(opts Options) (Executor, error)

type Options struct {
	IsStream    bool
	IsAnthropic bool
	Retry       int // 重试次数，0 表示不重试
}

type Registry struct {
	mu sync.RWMutex
	m  map[string]NewExecutorFunc
}

func (r *Registry) Register(provider string, fn NewExecutorFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.m[provider]; ok {
		panic(fmt.Errorf("provider %s already registered", provider))
	}
	r.m[provider] = fn
}

func (r *Registry) Get(provider string, opts Options) (Executor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fn, ok := r.m[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", provider)
	}
	return fn(opts)
}

var ExecutorRegistry = &Registry{
	m: make(map[string]NewExecutorFunc),
}
