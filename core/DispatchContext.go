package core

type DispatchContext struct {
	Dispatch map[uint64]bool
	Selector map[uint64][]string
	Result   map[uint64]Object
	Params   map[uint64][]Object
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
