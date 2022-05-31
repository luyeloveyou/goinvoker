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
	AddIfAbsent() (IFunctionChain, bool)
}

type ILibTable interface {
	core.ICoordinator
	Add(lib string, table IFunctionTable) bool
	AddIfAbsent(lib string) (IFunctionTable, bool)
}

type IInvokerTable interface {
	core.ICoordinator
	Add(invokerName string, table ILibTable) bool
	AddIfAbsent(invokerName string) (ILibTable, bool)
	SetUp(invokerName, libName, funcName, version string, handler ...func(reqId uint64, result any, params []any) (any, error)) bool
	SetUpAppend(invokerName, libName, funcName, version string, handler ...func(reqId uint64, result any, params []any) (any, error)) bool
	Call(invokerName, libName, funcName, version string, params ...any) (any, bool, error)
}
