package core

type DispatchContext struct {
	Dispatch map[uint64]bool
	Selector map[uint64][]string
	Result   map[uint64]any
	Params   map[uint64][]any
}

func NewDispatchContext() *DispatchContext {
	return &DispatchContext{
		Dispatch: make(map[uint64]bool),
		Selector: make(map[uint64][]string),
		Result:   make(map[uint64]any),
		Params:   make(map[uint64][]any),
	}
}

func (d *DispatchContext) clear(id uint64) {
	delete(d.Dispatch, id)
	delete(d.Selector, id)
	delete(d.Result, id)
	delete(d.Result, id)
	delete(d.Params, id)
}

func (d *DispatchContext) getDispatch(id uint64) bool {
	v, ok := d.Dispatch[id]
	if ok {
		return v
	} else {
		return false
	}
}

func (d *DispatchContext) setDispatch(id uint64, dispatch bool) {
	d.Dispatch[id] = dispatch
}

func (d *DispatchContext) getSelector(id uint64) []string {
	v, ok := d.Selector[id]
	if ok {
		return v
	} else {
		return nil
	}
}

func (d *DispatchContext) setSelector(id uint64, selector []string) {
	d.Selector[id] = selector
}

func (d *DispatchContext) getResult(id uint64) any {
	v, ok := d.Result[id]
	if ok {
		return v
	} else {
		return nil
	}
}

func (d *DispatchContext) setResult(id uint64, result any) {
	d.Result[id] = result
}

func (d *DispatchContext) getParams(id uint64) []any {
	v, ok := d.Params[id]
	if ok {
		return v
	} else {
		return nil
	}
}

func (d *DispatchContext) setParams(id uint64, params []any) {
	d.Params[id] = params
}
