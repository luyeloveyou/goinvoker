package functionchain

import (
	"goinvoker/chain"
)

type functionTable struct {
	RootChain chain.IFunctionChain
}

func NewFunctionTable() *functionTable {
	return &functionTable{}
}

func (f *functionTable) Chain() (chain.IFunctionChain, bool) {
	if f.RootChain == nil {
		return nil, false
	}
	return f.RootChain, true
}
