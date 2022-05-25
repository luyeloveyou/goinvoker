package core

type IHandler interface {
	IRouted
	Invoke(reqId uint64, result Object, params []Object) Object
}

type Handler struct {
	Routed
	HandleFunc func(reqId uint64, result Object, params []Object) Object
}

func NewHandler(handleFunc func(reqId uint64, result Object, params []Object) Object) *Handler {
	return &Handler{HandleFunc: handleFunc}
}

func (h *Handler) Invoke(reqId uint64, result Object, params []Object) Object {
	return h.HandleFunc(reqId, result, params)
}
