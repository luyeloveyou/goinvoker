package router

type FuncRouter struct {
	RouteFunc func(selector string) any
	AddFunc   func(selector string, routed any)
}

func NewFuncRouter(route func(selector string) any, add func(selector string, routed any)) *FuncRouter {
	return &FuncRouter{
		RouteFunc: route,
		AddFunc:   add,
	}
}

func (f *FuncRouter) Route(selector string) any {
	if f.RouteFunc == nil {
		return nil
	}
	return f.RouteFunc(selector)
}

func (f *FuncRouter) Add(selector string, routed any) {
	if f.AddFunc == nil {
		return
	}
	f.AddFunc(selector, routed)
}
