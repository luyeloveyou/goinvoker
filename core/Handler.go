package core

type IHandler interface {
	Next() Object
	Invoke(reqId uint64, result *Object, params []*Object) *Object
}

type Handler struct {
	NextRouted *Object
	HandleFunc func(reqId uint64, result *Object, params []*Object) *Object
}

func (h *Handler) Next() Object {
	return h.NextRouted
}

func (h *Handler) Invoke(reqId uint64, result *Object, params []*Object) *Object {
	return h.HandleFunc(reqId, result, params)
}
