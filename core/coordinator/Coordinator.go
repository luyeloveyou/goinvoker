package coordinator

import (
	"fmt"
	"github.com/luyeloveyou/goinvoker/core"
	"github.com/luyeloveyou/goinvoker/core/context"
	"github.com/luyeloveyou/goinvoker/core/routed"
	"math/rand"
	"time"
)

var dispatchContext = context.NewDispatchContext()

type Coordinator struct {
	*routed.Routed
	RootRouted any
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		Routed:     &routed.Routed{},
		RootRouted: nil,
	}
}

func (c *Coordinator) CanDispatch() bool {
	_, flag := c.Next()
	return flag
}

func (c *Coordinator) Invoke(reqId uint64, selectors []string, result any, params []any) (retResult any, retInvoked bool, retErr error) {
	if c.RootRouted == nil {
		retResult = nil
		retInvoked = false
		retErr = fmt.Errorf("当前函数簇为nil")
		return
	}
	if reqId == 0 {
		rand.Seed(time.Now().Unix())
		reqId = rand.Uint64()
		defer dispatchContext.Clear(reqId)
	}
	dispatchContext.SetDispatch(reqId, true)
	retResult, retInvoked, retErr = helper(c.RootRouted, reqId, selectors, result, params)
	if retErr != nil {
		return
	}
	if dispatchContext.GetDispatch(reqId) && c.CanDispatch() || !retInvoked {
		if dispatchContext.GetResult(reqId) != nil {
			retResult = dispatchContext.GetResult(reqId)
		}
		s := dispatchContext.GetSelector(reqId)
		if s == nil {
			s = selectors
		}
		p := dispatchContext.GetParams(reqId)
		if p == nil {
			p = params
		}
		next, _ := c.Next()
		tempResult, isInvoked, err := helper(next, reqId, s, retResult, p)
		if !retInvoked {
			retInvoked = isInvoked
		}
		if err != nil {
			retErr = err
			return
		}
		retResult = tempResult
	}
	return
}

func Dispatch(reqId uint64, selectors []string, result any, params []any) {
	dispatchContext.SetDispatch(reqId, true)
	dispatchContext.SetSelector(reqId, selectors)
	dispatchContext.SetParams(reqId, params)
	dispatchContext.SetResult(reqId, result)
}

func Block(reqId uint64) {
	dispatchContext.Clear(reqId)
}

func helper(routed any, reqId uint64, selectors []string, result any, params []any) (retResult any, retInvoked bool, retErr error) {
	retResult = result
	index := 0
	var ok bool
	var err error
	var tempResult any
	for routed != nil {
		switch r := routed.(type) {
		case core.IRouter:
			if index < len(selectors) {
				routed, ok = r.Route(selectors[index])
				if !ok {
					return
				}
			} else {
				retErr = fmt.Errorf("选择子数量 %d 错误", len(selectors))
				return
			}
			if index+1 < len(selectors) {
				index++
			}
		case core.IHandler:
			tempResult, err = r.Invoke(reqId, retResult, params)
			if err != nil {
				retErr = err
				return
			}
			retResult = tempResult
			retInvoked = true
			index = 0
			routed, ok = r.Next()
			if !ok {
				return
			}
		case core.ICoordinator:
			subSelectors := selectors[index:]
			tempResult, ok, err = r.Invoke(reqId, subSelectors, retResult, params)
			if !retInvoked {
				retInvoked = ok
			}
			if err != nil {
				retErr = err
				return
			}
			retResult = tempResult
			return
		default:
			retErr = fmt.Errorf("错误的类型: %T", routed)
			return
		}
	}
	return
}
