package table

import (
	"fmt"
	"goinvoker/chain"
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

func (i *InvokerTable) Add(invokerName string, table chain.ILibTable) {
	if i.RootRouted == nil {
		panic("根路由不能为空")
	}
	nameRouter, ok := i.RootRouted.(core.IRouter)
	if !ok {
		panic(fmt.Sprintf("根路由不是路由类型: %T", i.RootRouted))
	}
	nameRouter.Add(invokerName, table)
}

func (i *InvokerTable) Call(invokerName, libName, funcName, version string, params ...any) any {
	return i.Invoke(0, []string{invokerName, libName, funcName, version}, nil, params)
}
