package core

type IRouted interface {
	Next() (any, bool)
}

type IRouter interface {
	Route(selector string) (any, bool)
	Add(selector string, routed any) bool
}

type ICoordinator interface {
	IRouted
	CanDispatch() bool
	Invoke(reqId uint64, selectors []string, result any, params []any) (any, bool, error)
}

type IHandler interface {
	IRouted
	Invoke(reqId uint64, result any, params []any) (any, error)
}
