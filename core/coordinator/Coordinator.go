package coordinator

import (
	"fmt"
	"goinvoker/core"
	"goinvoker/core/context"
	"goinvoker/core/routed"
	"math/rand"
	"time"
)

type Coordinator struct {
	*routed.Routed
	RootRouted any
	Context    *context.DispatchContext
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		Routed:     &routed.Routed{},
		RootRouted: nil,
		Context:    context.NewDispatchContext(),
	}
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
	if c.Context.GetDispatch(reqId) && c.CanDispatch() {
		if c.Context.GetResult(reqId) != nil {
			tempResult = c.Context.GetResult(reqId)
		}
		s := c.Context.GetSelector(reqId)
		if s == nil {
			s = selectors
		}
		p := c.Context.GetParams(reqId)
		c.Context.Clear(reqId)
		tempResult = helper(c.Next(), reqId, s, tempResult, p)
	}
	return tempResult
}

func (c *Coordinator) Dispatch(reqId uint64, selectors []string, result any, params []any) {
	c.Context.SetDispatch(reqId, true)
	c.Context.SetSelector(reqId, selectors)
	c.Context.SetParams(reqId, params)
	c.Context.SetResult(reqId, result)
}

func helper(routed any, reqId uint64, selectors []string, result any, params []any) any {
	tempResult := result
	index := 0
	for routed != nil {
		switch r := routed.(type) {
		case core.IRouter:
			if index < len(selectors) {
				routed = r.Route(selectors[index])
			} else {
				panic(fmt.Sprintf("选择子数量 %d 错误", len(selectors)))
			}
			if index+1 < len(selectors) {
				index++
			}
		case core.IHandler:
			tempResult = r.Invoke(reqId, tempResult, params)
			index = 0
			routed = r.Next()
		case core.ICoordinator:
			subSelectors := selectors[index:]
			tempResult = r.Invoke(reqId, subSelectors, tempResult, params)
			routed = nil
		default:
			panic(fmt.Sprintf("错误的类型: %T", routed))
		}
	}
	return tempResult
}
