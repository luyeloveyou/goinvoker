package main

import (
	"fmt"
	"goinvoker/core"
)

func main() {
	system := produce("test")
	result := system.Invoke(0, []string{"test", "1.1.0"}, nil, nil)
	fmt.Println(result)
}

func produce(name string) *core.Coordinator {
	root := core.NewRouter()
	d1 := core.NewVersionRouter()
	d2 := core.NewRouter()
	d3 := core.NewVersionRouter()
	s1 := &core.Handler{
		HandleFunc: func(reqId uint64, result core.Object, params []core.Object) core.Object {
			v := "hello"
			fmt.Println("--------")
			fmt.Println(v)
			fmt.Println("--------")
			return v
		},
	}
	s2 := &core.Handler{HandleFunc: func(reqId uint64, result core.Object, params []core.Object) core.Object {
		v := fmt.Sprintf("%v world", result)
		fmt.Println("--------")
		fmt.Println(v)
		fmt.Println("--------")
		return v
	}}
	system := &core.Coordinator{
		RootRouted: root,
		Context: core.DispatchContext{
			Dispatch: make(map[uint64]bool),
			Selector: make(map[uint64][]string),
			Result:   make(map[uint64]core.Object),
			Params:   make(map[uint64][]core.Object),
		},
	}
	root.Add(name, d1)
	d1.Add("0.0.2", s1)
	s1.NextRouted = d2
	d2.Add(name, d3)
	d3.Add("1.0.0", s2)
	return system
}
