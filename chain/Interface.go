package chain

import (
	"goinvoker/core"
)

type IFunctionChain interface {
	core.ICoordinator
	Add(funcName, version string, handler core.IHandler)
}

type IFunctionTable interface {
	Chain() IFunctionChain
}

type ILibTable interface {
	core.ICoordinator
	Add(lib string, table IFunctionTable)
}

type IInvokerTable interface {
	core.ICoordinator
	Add(invokerName string, table ILibTable)
	Call(invokerName, libName, funcName, version string, params ...any) any
}
