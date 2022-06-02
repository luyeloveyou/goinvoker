package builder

import "github.com/luyeloveyou/goinvoker/core/coordinator"

type ChainCaller struct {
	c        *coordinator.Coordinator
	selector []string
	result   any
	params   []any
}

func NewChainCaller(c *coordinator.Coordinator) *ChainCaller {
	return &ChainCaller{c: c}
}

func (c *ChainCaller) Path(path ...string) *ChainCaller {
	c.selector = path
	return c
}

func (c *ChainCaller) Result(result any) *ChainCaller {
	c.result = result
	return c
}

func (c *ChainCaller) Params(params ...any) *ChainCaller {
	c.params = params
	return c
}

func (c *ChainCaller) Call() (any, bool, error) {
	s := c.selector
	r := c.result
	p := c.params
	c.selector = nil
	c.result = nil
	c.params = nil
	return c.c.Invoke(0, s, r, p)
}
