package chain

import (
	"goinvoker/core"
)

type IFunctionChain interface {
	core.ICoordinator
	Add(funcName, version string, handler core.IHandler) bool
}

type IFunctionTable interface {
	Chain() (IFunctionChain, bool)
}

type ILibTable interface {
	core.ICoordinator
	Add(lib string, table IFunctionTable) bool
}

type IInvokerTable interface {
	core.ICoordinator
	Add(invokerName string, table ILibTable) bool
	SetUp(invokerName, libName, funcName, version string, handler ...core.IHandler) bool
	Call(invokerName, libName, funcName, version string, params ...any) (any, bool, error)
}
