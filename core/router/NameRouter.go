package router

type NameRouter struct {
	cache map[string]any
}

func NewNameRouter() *NameRouter {
	return &NameRouter{cache: make(map[string]any)}
}

func (r *NameRouter) Route(selector string) any {
	return r.cache[selector]
}

func (r *NameRouter) Add(selector string, routed any) {
	r.cache[selector] = routed
}
