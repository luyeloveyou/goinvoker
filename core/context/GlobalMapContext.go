package context

import "sync"

type GlobalMapContext struct {
	cache *sync.Map
}

func NewGlobalMapContext() *GlobalMapContext {
	return &GlobalMapContext{&sync.Map{}}
}

func (g *GlobalMapContext) Clear(id uint64) {
	g.cache.Delete(id)
}

func (g *GlobalMapContext) Get(id uint64, key any) (any, bool) {
	value, ok := g.cache.Load(id)
	if ok {
		valueMap, ok := value.(map[any]any)
		if ok {
			v, ok := valueMap[key]
			if ok {
				return v, true
			}
		}
	}
	return nil, false
}

func (g *GlobalMapContext) Set(id uint64, key any, value any) bool {
	v, _ := g.cache.LoadOrStore(id, make(map[any]any))
	valueMap, ok := v.(map[any]any)
	if ok {
		valueMap[key] = value
		return true
	}
	return false
}
