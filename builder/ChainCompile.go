package builder

import (
	"fmt"
	"strings"
)

type ChainCompile struct {
	template  *ChainTemplate
	simpleDes []rune
	postfix   []rune
	index     int
	change    bool
	element   []any
}

func NewChainCompile(template *ChainTemplate) *ChainCompile {
	ret := &ChainCompile{}
	ret.template = template
	end := strings.LastIndexFunc(string(template.Description()), func(r rune) bool {
		return r == Infinite
	})
	ret.simpleDes = make([]rune, len(template.StructDes), cap(template.StructDes))
	copy(ret.simpleDes, template.StructDes)
	if end > 0 {
		temp := template.StructDes[0:end]
		start := strings.LastIndexFunc(string(temp), func(r rune) bool {
			return r == LP
		})
		ret.postfix = make([]rune, end-start-2)
		ret.simpleDes = make([]rune, len(template.StructDes)-end+start-1)
		i1 := 0
		i2 := 0
		for i := 0; i < len(template.StructDes); i++ {
			if i > start && i < end-1 {
				ret.postfix[i1] = template.StructDes[i]
				i1++
			} else if i < start || i > end {
				ret.simpleDes[i2] = template.StructDes[i]
				i2++
			}
		}
	}
	ret.simpleDes = cut(ret.simpleDes, func(r rune) bool {
		return r == LP || r == RP || r == Infinite
	})
	ret.index = 0
	ret.change = false
	return ret
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
		var token rune
		if i < len(compile.simpleDes) {
			token = compile.simpleDes[i]
		} else {
			token = compile.postfix[(i-len(compile.simpleDes))%len(compile.postfix)]
		}
		c.helper(e, token)
	}
	return c
}

func (c *ChainCompile) Resolve(compile *ChainCompile) *ChainCompile {
	c.Clear()
	for i, e := range compile.element {
		var token rune
		if i < len(compile.simpleDes) {
			token = compile.simpleDes[i]
		} else {
			token = compile.postfix[(i-len(compile.simpleDes))%len(compile.postfix)]
		}
		c.helper(e, token)
	}
	return c
}

func (c *ChainCompile) Clear() *ChainCompile {
	c.element = []any{}
	c.index = 0
	c.change = false
	return c
}

func (c *ChainCompile) Compile() *ChainCompile {
	if len(c.postfix) == 0 && len(c.element) != len(c.simpleDes) || len(c.postfix) != 0 && (len(c.element)-len(c.simpleDes))%len(c.postfix) != 0 {
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
	if c.index >= len(c.simpleDes) && !c.change {
		if len(c.postfix) != 0 {
			c.change = true
			c.index = 0
		} else {
			panic("参数超过长度限制")
		}
	}
	if c.change && c.postfix[c.index%len(c.postfix)] == token || !c.change && c.simpleDes[c.index] == token {
		c.index++
		c.element = append(c.element, value)
	} else {
		if c.change {
			panic(fmt.Sprintf("当前需求为:%v, 传入类型为:%v", string(c.postfix[c.index%len(c.postfix)]), string(token)))
		} else {
			panic(fmt.Sprintf("当前需求为:%v, 传入类型为:%v", string(c.simpleDes[c.index]), string(token)))
		}
	}
	return c
}

func cut(s []rune, f func(rune) bool) []rune {
	tmp := make([]rune, 0, 0)
	for _, r := range s {
		if !f(r) {
			tmp = append(tmp, r)
		}
	}
	return tmp
}
