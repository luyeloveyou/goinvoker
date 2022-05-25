package core

import (
	"fmt"
	"math/rand"
	"time"
)

type ICoordinator interface {
	Root() Object
	Next() Object
	CanDispatch() bool
	Invoke(reqId uint64, selectors []string, result *Object, params []*Object) *Object
	Dispatch(reqId uint64, selectors []string, result *Object, params []*Object) *Object
}

type Coordinator struct {
	RootRouted *Object
	NextRouted *Object
	Context    DispatchContext
}

func (c *Coordinator) Next() Object {
	return c.NextRouted
}

func (c *Coordinator) Root() Object {
	return c.RootRouted
}

func (c *Coordinator) CanDispatch() bool {
	return c.Next() != nil
}

func (c *Coordinator) Invoke(reqId uint64, selectors []string, result *Object, params []*Object) *Object {
	if c.RootRouted == nil {
		panic("当前函数簇为nil")
	}
	rand.Seed(time.Now().Unix())
	if reqId == 0 {
		reqId = rand.Uint64()
	}
	tempResult := helper(c.RootRouted, reqId, selectors, result, params)
	return tempResult
}

func (c *Coordinator) Dispatch(reqId uint64, selectors []string, result *Object, params []*Object) *Object {
	//TODO implement me
	panic("implement me")
}

func helper(routed Object, reqId uint64, selectors []string, result *Object, params []*Object) *Object {
	tempResult := result
	index := 0
	for routed != nil {
		switch routed.(type) {
		case IRouter:
			var router = routed.(IRouter)
			if index < len(selectors) {
				routed = router.Route(selectors[index])
			} else {
				panic(fmt.Sprintf("选择子数量 %d 错误", len(selectors)))
			}
			if index+1 < len(selectors) {
				index++
			}
		case IHandler:
			var handler = routed.(IHandler)
			tempResult = handler.Invoke(reqId, tempResult, params)
			index = 0
			routed = handler.Next()
		case ICoordinator:
			var coordinator = routed.(ICoordinator)
			subSelectors := selectors[index:]
			tempResult = coordinator.Invoke(reqId, subSelectors, tempResult, params)
		default:
			panic(fmt.Sprintf("错误的类型: %T", routed))
		}
	}
	return tempResult
}
