package main

import (
	"fmt"
	"goinvoker/core/coordinator"
	"goinvoker/core/handler"
	"goinvoker/core/router"
)

func main() {
	system := produce("test")
	system1 := produce("test")
	system.NextRouted = system1
	system.Dispatch(1234, []string{"test", "1.1.0"}, nil, nil)
	result := system.Invoke(1234, []string{"test", "1.1.0"}, nil, nil)
	fmt.Println(result)
	//a := &A{Aname: "a"}
	//var b = &B{
	//	A:     a,
	//	Bname: "b",
	//}
	//a.Aname = "test"
	//fmt.Printf("%v\n", b.A)
	//fmt.Println(a)
}

type IA interface {
	AN() string
}

type A struct {
	Aname string
}

type IB interface {
	IA
	BN() string
}

type B struct {
	*A
	Bname string
}

func (a *A) AN() string {
	a.Aname = "AN"
	return a.Aname
}

func (b *B) BN() string {
	b.Bname = "BN"
	b.Aname = "BN"
	return b.Bname
}

func test(b *B) {
	b.Aname = "testA"
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
