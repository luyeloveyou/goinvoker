package router

type FuncRouter struct {
	RouteFunc func(selector string) (any, bool)
	AddFunc   func(selector string, routed any) bool
}

func NewFuncRouter(routeFunc func(selector string) (any, bool), addFunc func(selector string, routed any) bool) *FuncRouter {
	return &FuncRouter{RouteFunc: routeFunc, AddFunc: addFunc}
}

func (f *FuncRouter) Route(selector string) (any, bool) {
	if f.RouteFunc == nil {
		return nil, false
	}
	return f.RouteFunc(selector)
}

func (f *FuncRouter) Add(selector string, routed any) bool {
	if f.AddFunc == nil {
		return false
	}
	return f.AddFunc(selector, routed)
}

func (f *FuncRouter) Has(selector string) (any, bool) {
	return f.Route(selector)
}
