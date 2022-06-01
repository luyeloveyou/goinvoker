package builder

import (
	"fmt"
	"testing"
)

func TestTemplate(t *testing.T) {
	functionChain := NewChainTemplate().Name().Version().Handle()
	invokerTable := NewChainTemplate().Name().Name().Next(functionChain).Next(functionChain)
	functionChain.Fill("func", "0.0.0")
	clone := functionChain.Clone()
	functionChain.Append(HandlerHelper(func(id uint64, result any, params []any) (any, error) {
		return 1, nil
	}))
	clone.Append(HandlerHelper(func(id uint64, result any, params []any) (any, error) {
		return 2, nil
	}))
	invokerTable.Fill("invoker", "lib", functionChain, clone)
	invokerTable.Build()
	invoke, _, err := invokerTable.Coordinator().Invoke(0, []string{"invoker", "lib", "func", "0.0.0"}, nil, nil)
	if err != nil {
		return
	}
	fmt.Println(invoke)
}
