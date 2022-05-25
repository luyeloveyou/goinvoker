package core

import (
	"fmt"
	"math/rand"
	"time"
)

type ICoordinator interface {
	IRouted
	CanDispatch() bool
	Invoke(reqId uint64, selectors []string, result any, params []any) any
	Dispatch(reqId uint64, selectors []string, result any, params []any)
}

type Coordinator struct {
	Routed
	RootRouted any
	Context    *DispatchContext
}

func NewCoordinator() *Coordinator {
	return &Coordinator{Context: NewDispatchContext()}
}

func (c *Coordinator) CanDispatch() bool {
	return c.Next() != nil
}

func (c *Coordinator) Invoke(reqId uint64, selectors []string, result any, params []any) any {
	if c.RootRouted == nil {
		panic("当前函数簇为nil")
	}
	if reqId == 0 {
		rand.Seed(time.Now().Unix())
		reqId = rand.Uint64()
	}
	tempResult := helper(c.RootRouted, reqId, selectors, result, params)
	if c.Context.getDispatch(reqId) && c.CanDispatch() {
		if c.Context.getResult(reqId) != nil {
			tempResult = c.Context.getResult(reqId)
		}
		s := c.Context.getSelector(reqId)
		p := c.Context.getParams(reqId)
		c.Context.clear(reqId)
		tempResult = helper(c.Next(), reqId, s, tempResult, p)
	}
	return tempResult
}

func (c *Coordinator) Dispatch(reqId uint64, selectors []string, result any, params []any) {
	c.Context.setDispatch(reqId, true)
	c.Context.setSelector(reqId, selectors)
	c.Context.setParams(reqId, params)
	c.Context.setResult(reqId, result)
}

func helper(routed any, reqId uint64, selectors []string, result any, params []any) any {
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
			routed = nil
		default:
			panic(fmt.Sprintf("错误的类型: %T", routed))
		}
	}
	return tempResult
}
