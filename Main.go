package main

import (
	"fmt"
	"goinvoker/chain/functionchain"
	"goinvoker/chain/table"
	"goinvoker/core/coordinator"
	"goinvoker/core/handler"
	"goinvoker/core/router"
)

func main() {
	validChain := functionchain.NewFunctionChain()
	addChain := functionchain.NewFunctionChain()
	addTable := functionchain.NewFunctionTable()
	addLib := table.NewLibTable()
	addInvoker := table.NewInvokerTable()
	validChain.Add("add", "1.0.0", handler.NewHandler(func(reqId uint64, result any, params []any) any {
		if valid(params[0].(int), params[1].(int)) {
			validChain.Dispatch(reqId, nil, result, params)
		} else {
			fmt.Println("error")
		}
		return 0
	}))
	addChain.Add("add", "0.0.0", handler.NewHandler(func(reqId uint64, result any, params []any) any {
		return add(params[0].(int), params[1].(int))
	}))
	validChain.NextRouted = addChain
	addTable.RootChain = validChain
	addLib.Add("add", addTable)
	addInvoker.Add("add", addLib)
	call := addInvoker.Call("add", "add", "add", "1.1.0", 1, -2)
	fmt.Println(call)
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

func produce(name string) *coordinator.Coordinator {
	root := router.NewNameRouter()
	d1 := router.NewVersionRouter()
	d2 := router.NewNameRouter()
	d3 := router.NewVersionRouter()
	s1 := handler.NewHandler(func(reqId uint64, result any, params []any) any {
		var v string
		if result != nil {
			v = fmt.Sprintf("hello %v", result)
		} else {
			v = "hello"
		}
		fmt.Println("--------")
		fmt.Println(v)
		fmt.Println("--------")
		return v
	})
	s2 := handler.NewHandler(func(reqId uint64, result any, params []any) any {
		v := fmt.Sprintf("%v world", result)
		fmt.Println("--------")
		fmt.Println(v)
		fmt.Println("--------")
		return v
	})
	system := coordinator.NewCoordinator()
	system.RootRouted = root
	root.Add(name, d1)
	d1.Add("0.0.2", s1)
	s1.NextRouted = d2
	d2.Add(name, d3)
	d3.Add("1.0.0", s2)
	return system
}
