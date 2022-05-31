package table

import (
	"github.com/luyeloveyou/goinvoker/chain"
	"github.com/luyeloveyou/goinvoker/chain/functionchain"
	"github.com/luyeloveyou/goinvoker/core/coordinator"
	"github.com/luyeloveyou/goinvoker/core/router"
)

type LibTable struct {
	*coordinator.Coordinator
	Tables map[string]chain.IFunctionTable
}

func NewLibTable() *LibTable {
	libTable := &LibTable{
		Coordinator: coordinator.NewCoordinator(),
		Tables:      make(map[string]chain.IFunctionTable),
	}
	libTable.RootRouted = router.NewFuncRouter(func(selector string) (any, bool) {
		v, ok := libTable.Tables[selector]
		if !ok {
			return nil, false
		}
		functionTable, ok := v.(chain.IFunctionTable)
		if !ok {
			return nil, false
		}
		fc, b := functionTable.Chain()
		return fc, b
	}, nil)
	return libTable
}

func (t *LibTable) Add(lib string, table chain.IFunctionTable) bool {
	if t.Tables == nil {
		return false
	}
	t.Tables[lib] = table
	return true
}

func (t *LibTable) AddIfAbsent(lib string) (chain.IFunctionTable, bool) {
	if t.Tables == nil {
		return nil, false
	}
	ft, ok := t.Tables[lib]
	if !ok {
		ft = functionchain.NewFunctionTable()
		t.Tables[lib] = ft
	}
	return ft, true
}
