package context

var (
	dispatch = "DISPATCH"
	selector = "SELECTOR"
	result   = "RESULT"
	params   = "PARAMS"
)

type DispatchContext struct {
	cache *GlobalMapContext
}

func NewDispatchContext() *DispatchContext {
	return &DispatchContext{
		NewGlobalMapContext(),
	}
}

func (d *DispatchContext) Clear(id uint64) {
	d.cache.Clear(id)
}

func get(d *DispatchContext, id uint64, key any) any {
	v, ok := d.cache.Get(id, key)
	if ok {
		return v
	}
	return nil
}

func (d *DispatchContext) GetDispatch(id uint64) bool {
	v := get(d, id, dispatch)
	if v == nil {
		return false
	}
	return v.(bool)
}

func (d *DispatchContext) SetDispatch(id uint64, value bool) {
	d.cache.Set(id, dispatch, value)
}

func (d *DispatchContext) GetSelector(id uint64) []string {
	v := get(d, id, selector)
	if v == nil {
		return nil
	}
	return v.([]string)
}

func (d *DispatchContext) SetSelector(id uint64, value []string) {
	d.cache.Set(id, selector, value)
}

func (d *DispatchContext) GetResult(id uint64) any {
	return get(d, id, result)
}

func (d *DispatchContext) SetResult(id uint64, value any) {
	d.cache.Set(id, result, value)
}

func (d *DispatchContext) GetParams(id uint64) []any {
	v := get(d, id, params)
	if v == nil {
		return nil
	}
	return v.([]any)
}

func (d *DispatchContext) SetParams(id uint64, value []any) {
	d.cache.Set(id, params, value)
}
