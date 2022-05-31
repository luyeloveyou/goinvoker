package functionchain

import (
	"github.com/luyeloveyou/goinvoker/chain"
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

func (f *FunctionTable) AddIfAbsent() (chain.IFunctionChain, bool) {
	if f.RootChain == nil {
		f.RootChain = NewFunctionChain()
	}
	return f.RootChain, true
}
