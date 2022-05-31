package table

import (
	"goinvoker/chain"
	"goinvoker/core/coordinator"
	"goinvoker/core/router"
)

type libTable struct {
	*coordinator.Coordinator
	Tables map[string]chain.IFunctionTable
}

func NewLibTable() *libTable {
	libTable := &libTable{
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

func (t *libTable) Add(lib string, table chain.IFunctionTable) bool {
	if t.Tables == nil {
		return false
	}
	t.Tables[lib] = table
	return true
}
