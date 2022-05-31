package table

import (
	"github.com/luyeloveyou/goinvoker/chain"
	"github.com/luyeloveyou/goinvoker/chain/functionchain"
	"github.com/luyeloveyou/goinvoker/core"
	"github.com/luyeloveyou/goinvoker/core/coordinator"
	"github.com/luyeloveyou/goinvoker/core/handler"
	"github.com/luyeloveyou/goinvoker/core/router"
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

func (i *InvokerTable) AddIfAbsent(invokerName string) (chain.ILibTable, bool) {
	if i.RootRouted == nil {
		return nil, false
	}
	nameRouter, ok := i.RootRouted.(core.IRouter)
	if !ok {
		return nil, false
	}
	temp, ok := nameRouter.Route(invokerName)
	if !ok {
		temp = NewLibTable()
		nameRouter.Add(invokerName, temp)
	}
	lt := temp.(chain.ILibTable)
	return lt, true
}

func (i *InvokerTable) SetUpAppend(invokerName, libName, funcName, version string, hs ...func(reqId uint64, result any, params []any) (any, error)) bool {
	temp, ok := i.AddIfAbsent(invokerName)
	if !ok {
		return false
	}
	lt := temp.(*LibTable)

	ft, ok := lt.AddIfAbsent(libName)
	if !ok {
		return false
	}

	fc, ok := ft.AddIfAbsent()
	if !ok {
		return false
	}

	header := fc.(*functionchain.FunctionChain)
	for header.NextRouted != nil {
		header = header.NextRouted.(*functionchain.FunctionChain)
	}
	if len(hs) > 0 {
		header.NextRouted = functionchain.NewFunctionChain()
		header = header.NextRouted.(*functionchain.FunctionChain)
	}
	for index := 0; index < len(hs)-1; index++ {
		header.Add(funcName, version, handler.NewHandler(hs[index]))
		if header.NextRouted == nil {
			header.NextRouted = functionchain.NewFunctionChain()
		}
		header = header.NextRouted.(*functionchain.FunctionChain)
	}
	if len(hs) > 0 {
		header.Add(funcName, version, handler.NewHandler(hs[len(hs)-1]))
	}
	return true
}

func (i *InvokerTable) SetUp(invokerName, libName, funcName, version string, hs ...func(reqId uint64, result any, params []any) (any, error)) bool {
	temp, ok := i.AddIfAbsent(invokerName)
	if !ok {
		return false
	}
	lt := temp.(*LibTable)

	ft, ok := lt.AddIfAbsent(libName)
	if !ok {
		return false
	}

	fc, ok := ft.AddIfAbsent()
	if !ok {
		return false
	}

	header := fc.(*functionchain.FunctionChain)
	for index := 0; index < len(hs)-1; index++ {
		header.Add(funcName, version, handler.NewHandler(hs[index]))
		if header.NextRouted == nil {
			header.NextRouted = functionchain.NewFunctionChain()
		}
		header = header.NextRouted.(*functionchain.FunctionChain)
	}
	if len(hs) > 0 {
		header.Add(funcName, version, handler.NewHandler(hs[len(hs)-1]))
	}
	return true
}

func (i *InvokerTable) Call(invokerName, libName, funcName, version string, params ...any) (any, bool, error) {
	return i.Invoke(0, []string{invokerName, libName, funcName, version}, nil, params)
}
