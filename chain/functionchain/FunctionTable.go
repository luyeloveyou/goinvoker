package functionchain

import (
	"goinvoker/chain"
)

type FunctionTable struct {
	RootChain chain.IFunctionChain
}

func NewFunctionTable() *FunctionTable {
	return &FunctionTable{}
}

func (f *FunctionTable) Chain() (chain.IFunctionChain, bool) {
	if f.RootChain == nil {
		return nil, false
	}
	return f.RootChain, true
}
