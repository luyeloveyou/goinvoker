package context

import (
	"context"
	"sync"
)

type MapContext struct {
	context.Context
	values *sync.Map
}

func (g *MapContext) Value(key any) any {
	v, ok := g.values.Load(key)
	if ok {
		return v
	} else {
		return g.Context.Value(key)
	}
}

func (g *MapContext) Put(key, value any) {
	g.values.Store(key, value)
}

func WithMap(ctx context.Context) *MapContext {
	return &MapContext{
		Context: ctx,
		values:  &sync.Map{},
	}
}
