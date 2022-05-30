package table

import (
	"goinvoker/chain"
	"goinvoker/core/coordinator"
	"goinvoker/core/router"
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
	libTable.RootRouted = router.NewFuncRouter(func(selector string) any {
		v, ok := libTable.Tables[selector]
		if !ok {
			return nil
		}
		functionTable, ok := v.(chain.IFunctionTable)
		if !ok {
			return nil
		}
		return functionTable.Chain()
	}, nil)
	return libTable
}

func (t *LibTable) Add(lib string, table chain.IFunctionTable) {
	if t.Tables == nil {
		panic("库表不能为nil")
	}
	t.Tables[lib] = table
}
