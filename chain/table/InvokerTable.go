package table

import (
	"goinvoker/chain"
	"goinvoker/core"
	"goinvoker/core/coordinator"
	"goinvoker/core/router"
)

type invokerTable struct {
	*coordinator.Coordinator
}

func NewInvokerTable() *invokerTable {
	invokerTable := &invokerTable{coordinator.NewCoordinator()}
	invokerTable.RootRouted = router.NewNameRouter()
	return invokerTable
}

func (i *invokerTable) Add(invokerName string, table chain.ILibTable) bool {
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

func (i *invokerTable) Call(invokerName, libName, funcName, version string, params ...any) (any, bool, error) {
	return i.Invoke(0, []string{invokerName, libName, funcName, version}, nil, params)
}
