package core

type IHandler interface {
	IRouted
	Invoke(reqId uint64, result any, params []any) any
}

type Handler struct {
	Routed
	HandleFunc func(reqId uint64, result any, params []any) any
}

func NewHandler(handleFunc func(reqId uint64, result any, params []any) any) *Handler {
	return &Handler{HandleFunc: handleFunc}
}

func (h *Handler) Invoke(reqId uint64, result any, params []any) any {
	return h.HandleFunc(reqId, result, params)
}
