package router

type NameRouter struct {
	cache map[string]any
}

func NewNameRouter() *NameRouter {
	return &NameRouter{cache: make(map[string]any)}
}

func (r *NameRouter) Route(selector string) (any, bool) {
	v, ok := r.cache[selector]
	return v, ok
}

func (r *NameRouter) Add(selector string, routed any) bool {
	r.cache[selector] = routed
	return true
}

func (r *NameRouter) Has(selector string) (any, bool) {
	return r.Route(selector)
}
