package handler

import (
	"goinvoker/core/routed"
)

type Handler struct {
	*routed.Routed
	HandleFunc func(reqId uint64, result any, params []any) any
}

func NewHandler(handleFunc func(reqId uint64, result any, params []any) any) *Handler {
	return &Handler{
		Routed:     &routed.Routed{},
		HandleFunc: handleFunc,
	}
}

func (h *Handler) Invoke(reqId uint64, result any, params []any) any {
	return h.HandleFunc(reqId, result, params)
}