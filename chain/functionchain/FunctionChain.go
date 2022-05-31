package functionchain

import (
	"goinvoker/core"
	"goinvoker/core/coordinator"
	"goinvoker/core/router"
)

type FunctionChain struct {
	*coordinator.Coordinator
}

func NewFunctionChain() *FunctionChain {
	functionChain := &FunctionChain{
		Coordinator: coordinator.NewCoordinator(),
	}
	functionChain.RootRouted = router.NewNameRouter()
	return functionChain
}

func (f *FunctionChain) Add(funcName, version string, handler core.IHandler) bool {
	if f.RootRouted == nil {
		return false
	}
	nameRouter, ok := f.RootRouted.(*router.NameRouter)
	if !ok {
		return false
	}
	routed, _ := nameRouter.Route(funcName)
	if routed == nil {
		routed = router.NewVersionRouter()
		nameRouter.Add(funcName, routed)
	}
	versionRouter, ok := routed.(*router.VersionRouter)
	if !ok {
		return false
	}
	versionRouter.Add(version, handler)
	return true
}

func DispatchIdRP(reqId uint64, result any, params ...any) {
	coordinator.Dispatch(reqId, nil, result, params)
}

func DispatchIdR(reqId uint64, result any) {
	coordinator.Dispatch(reqId, nil, result, nil)
}

func DispatchId(reqId uint64) {
	coordinator.Dispatch(reqId, nil, nil, nil)
}
