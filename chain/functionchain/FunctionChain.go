package functionchain

import (
	"goinvoker/core"
	"goinvoker/core/coordinator"
	"goinvoker/core/router"
)

type functionChain struct {
	*coordinator.Coordinator
}

func NewFunctionChain() *functionChain {
	functionChain := &functionChain{
		Coordinator: coordinator.NewCoordinator(),
	}
	functionChain.RootRouted = router.NewNameRouter()
	return functionChain
}

func (f *functionChain) Add(funcName, version string, handler core.IHandler) bool {
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

func DispatchIdRP(reqId uint64, result any, params []any) {
	coordinator.Dispatch(reqId, nil, result, params)
}

func DispatchIdR(reqId uint64, result any) {
	DispatchIdRP(reqId, result, nil)
}

func DispatchId(reqId uint64) {
	DispatchIdR(reqId, nil)
}
