package table

import (
	"goinvoker/chain"
	"goinvoker/chain/functionchain"
	"goinvoker/core"
	"goinvoker/core/coordinator"
	"goinvoker/core/router"
)

type InvokerTable struct {
	*coordinator.Coordinator
}

func NewInvokerTable() *InvokerTable {
	invokerTable := &InvokerTable{coordinator.NewCoordinator()}
	invokerTable.RootRouted = router.NewNameRouter()
	return invokerTable
}

func (i *InvokerTable) Add(invokerName string, table chain.ILibTable) bool {
	if i.RootRouted == nil {
		return false
	}
	nameRouter, ok := i.RootRouted.(core.IRouter)
	if !ok {
		return false
	}
	nameRouter.Add(invokerName, table)
	return true
}

func (i *InvokerTable) SetUp(invokerName, libName, funcName, version string, handler ...core.IHandler) bool {
	if i.RootRouted == nil {
		return false
	}
	nameRouter, ok := i.RootRouted.(core.IRouter)
	if !ok {
		return false
	}
	temp, ok := nameRouter.Route(invokerName)
	if !ok {
		temp = NewLibTable()
		nameRouter.Add(invokerName, temp)
	}
	lt := temp.(*LibTable)
	if lt.Tables == nil {
		return false
	}
	ft, ok := lt.Tables[libName]
	if !ok {
		ft = functionchain.NewFunctionTable()
		lt.Tables[libName] = ft
	}
	functionTable := ft.(*functionchain.FunctionTable)
	if functionTable.RootChain == nil {
		functionTable.RootChain = functionchain.NewFunctionChain()
	}
	header := functionTable.RootChain
	for index := 0; index < len(handler)-1; index++ {
		header.Add(funcName, version, handler[index])
		fc := header.(*functionchain.FunctionChain)
		if fc.NextRouted == nil {
			fc.NextRouted = functionchain.NewFunctionChain()
		}
		header = fc.NextRouted.(chain.IFunctionChain)
	}
	if len(handler) > 1 {
		header.Add(funcName, version, handler[len(handler)-1])
	}
	return true
}

func (i *InvokerTable) Call(invokerName, libName, funcName, version string, params ...any) (any, bool, error) {
	return i.Invoke(0, []string{invokerName, libName, funcName, version}, nil, params)
}
