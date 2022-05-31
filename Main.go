package main

import (
	"fmt"
	"goinvoker/chain/functionchain"
	"goinvoker/chain/table"
	"goinvoker/core/handler"
)

func main() {
	testFunctionChain()
}

func testFunctionChain() {
	validChain := functionchain.NewFunctionChain()
	addChain := functionchain.NewFunctionChain()
	addTable := functionchain.NewFunctionTable()
	addLib := table.NewLibTable()
	addInvoker := table.NewInvokerTable()
	validChain.Add("add", "1.0.0", handler.NewHandler(func(reqId uint64, result any, params []any) (any, error) {
		if !valid(params[0].(int), params[1].(int)) {
			fmt.Println("error")
		}
		functionchain.DispatchId(reqId)
		return 0, nil
	}))
	addChain.Add("add", "0.2.0", handler.NewHandler(func(reqId uint64, result any, params []any) (any, error) {
		return add(params[0].(int), params[1].(int)), nil
	}))
	validChain.NextRouted = addChain
	addTable.RootChain = validChain
	addLib.Add("add", addTable)
	addInvoker.Add("add", addLib)
	call, b, err := addInvoker.Call("add", "add", "add", "0.1.0", 1, -2)
	fmt.Println(call, b, err)
}

func valid(a, b int) bool {
	if a < 0 || b < 0 {
		return false
	} else {
		return true
	}
}

func add(a, b int) int {
	return a + b
}
