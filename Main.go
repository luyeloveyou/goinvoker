package main

import (
	"fmt"
	"goinvoker/core"
)

func main() {
	system := produce("test")
	system1 := produce("test")
	system.NextRouted = system1
	fmt.Println(system.Next())
	system.Dispatch(1234, []string{"test", "1.1.0"}, nil, nil)
	result := system.Invoke(1234, []string{"test", "1.1.0"}, nil, nil)
	fmt.Println(result)
}

func produce(name string) *core.Coordinator {
	root := core.NewRouter()
	d1 := core.NewVersionRouter()
	d2 := core.NewRouter()
	d3 := core.NewVersionRouter()
	s1 := core.NewHandler(func(reqId uint64, result core.Object, params []core.Object) core.Object {
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
	s2 := core.NewHandler(func(reqId uint64, result core.Object, params []core.Object) core.Object {
		v := fmt.Sprintf("%v world", result)
		fmt.Println("--------")
		fmt.Println(v)
		fmt.Println("--------")
		return v
	})
	system := core.NewCoordinator()
	system.RootRouted = root
	root.Add(name, d1)
	d1.Add("0.0.2", s1)
	s1.NextRouted = d2
	d2.Add(name, d3)
	d3.Add("1.0.0", s2)
	return system
}
