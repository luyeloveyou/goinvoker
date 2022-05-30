package core

type IRouted interface {
	Next() any
}

type IRouter interface {
	Route(selector string) any
	Add(selector string, routed any)
}

type ICoordinator interface {
	IRouted
	CanDispatch() bool
	Invoke(reqId uint64, selectors []string, result any, params []any) any
	Dispatch(reqId uint64, selectors []string, result any, params []any)
}

type IHandler interface {
	IRouted
	Invoke(reqId uint64, result any, params []any) any
}
