package builder

import (
	"fmt"
)

type ChainCompile struct {
	template  *ChainTemplate
	simpleDes []rune
	index     int
	element   []any
}

func NewChainCompile(template *ChainTemplate) *ChainCompile {
	return &ChainCompile{
		template:  template,
		simpleDes: cut(template.Description(), func(r rune) bool { return r == LP || r == RP }),
		index:     0,
		element:   []any{},
	}
}

func (c *ChainCompile) Name(name string) *ChainCompile {
	return c.helper(name, NameRouter)
}

func (c *ChainCompile) Version(version string) *ChainCompile {
	return c.helper(version, VersionRouter)
}

func (c *ChainCompile) Handle(handle func(id uint64, result any, params []any) (any, error)) *ChainCompile {
	return c.helper(handle, Handler)
}

func (c *ChainCompile) Next(compile *ChainCompile) *ChainCompile {
	for i, e := range compile.element {
		c.helper(e, compile.simpleDes[i])
	}
	return c
}

func (c *ChainCompile) Resolve(compile *ChainCompile) *ChainCompile {
	c.Clear()
	for i, e := range compile.element {
		c.helper(e, compile.simpleDes[i])
	}
	return c
}

func (c *ChainCompile) Clear() *ChainCompile {
	c.element = []any{}
	c.index = 0
	return c
}

func (c *ChainCompile) Compile() *ChainCompile {
	if len(c.element) != len(c.simpleDes) {
		panic("参数长度不一致")
	}
	c.template.Fill(c.element...)
	c.template.Build()
	return c.Clear()
}

func (c *ChainCompile) Result() *ChainCaller {
	return c.template.Result()
}

func (c *ChainCompile) helper(value any, token rune) *ChainCompile {
	if c.simpleDes[c.index] == token {
		c.index++
		c.element = append(c.element, value)
	} else {
		panic(fmt.Sprintf("当前需求为:%v, 传入类型为:%v", string(c.template.StructDes[c.index]), string(token)))
	}
	return c
}

func cut(s []rune, f func(rune) bool) []rune {
	tmp := make([]rune, 0, len(s)*4/5)
	for _, r := range s {
		if !f(r) {
			tmp = append(tmp, r)
		}
	}
	return tmp
}
