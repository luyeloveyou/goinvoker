package main

import (
	"fmt"
	"goinvoker/chain/functionchain"
	"goinvoker/chain/table"
)

func main() {
	testSetup()
}

func testSetup() {
	invokerTable := table.NewInvokerTable()
	invokerTable.SetUp("add", "addLib", "addFunc", "1.1.0", valid)

	invokerTable.SetUpAppend("add", "addLib", "addFunc", "1.1.0", add)
	call, b, err := invokerTable.Call("add", "addLib", "addFunc", "1.2.0", 1, 2)
	fmt.Println(call, b, err)
}

func valid(reqId uint64, result any, params []any) (any, error) {
	a := params[0].(int)
	b := params[1].(int)

	if a < 0 || b < 0 {
		fmt.Println("error")
	}
	functionchain.DispatchId(reqId)
	return 2, nil
}

func add(reqId uint64, result any, params []any) (any, error) {
	a := params[0].(int)
	b := params[1].(int)
	functionchain.DispatchIdRP(reqId, nil, 2*params[0].(int), 2*params[1].(int))
	return a + b + result.(int), nil
}
