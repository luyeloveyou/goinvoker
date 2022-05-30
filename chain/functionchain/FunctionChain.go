package functionchain

import (
	"fmt"
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
	return functionChain
}

func (f *FunctionChain) Add(funcName, version string, handler core.IHandler) {
	if f.RootRouted == nil {
		panic("根路由不能为nil")
	}
	nameRouter, ok := f.RootRouted.(*router.NameRouter)
	if !ok {
		panic(fmt.Sprintf("根路由器类型 %T 不是名称路由器", f.RootRouted))
	}
	routed := nameRouter.Route(funcName)
	if routed == nil {
		routed = router.NewVersionRouter()
		nameRouter.Add(funcName, routed)
	}
	versionRouter, ok := routed.(*router.VersionRouter)
	if !ok {
		panic(fmt.Sprintf("名称路由器的路由结果类型 %T 不是版本路由器", routed))
	}
	versionRouter.Add(version, handler)
}
