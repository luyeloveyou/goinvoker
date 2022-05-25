package core

type IStrategy interface {
	Handler() *IHandler
	Strategy(reqId uint64, result Object, params []Object) Object
}
