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

func (f *FunctionTable) Chain() chain.IFunctionChain {
	return f.RootChain
}
