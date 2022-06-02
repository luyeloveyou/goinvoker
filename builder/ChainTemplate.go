package builder

import (
	"github.com/luyeloveyou/goinvoker/core"
	"github.com/luyeloveyou/goinvoker/core/coordinator"
	"github.com/luyeloveyou/goinvoker/core/handler"
	"github.com/luyeloveyou/goinvoker/core/router"
)

const (
	NameRouter    = 'N'
	VersionRouter = 'V'
	Handler       = 'H'
	LP            = '('
	RP            = ')'
)

type ChainTemplate struct {
	StructDes []rune
	Element   []any
	nodeH     []*coordinator.Coordinator
	lock      bool
}

func NewChainTemplate() *ChainTemplate {
	return &ChainTemplate{
		StructDes: []rune{},
		Element:   []any{},
		nodeH:     make([]*coordinator.Coordinator, 0, 4),
		lock:      false,
	}
}

func (c *ChainTemplate) Name() *ChainTemplate {
	if c.lock {
		panic("已生成的模板不可修改")
	}
	c.StructDes = append(c.StructDes, NameRouter)
	return c
}

func (c *ChainTemplate) Version() *ChainTemplate {
	if c.lock {
		panic("已生成的模板不可修改")
	}
	c.StructDes = append(c.StructDes, VersionRouter)
	return c
}

func (c *ChainTemplate) Next(template *ChainTemplate) *ChainTemplate {
	if c.lock {
		panic("已生成的模板不可修改")
	}
	description := template.Description()
	c.StructDes = append(c.StructDes, LP)
	c.StructDes = append(c.StructDes, description...)
	c.StructDes = append(c.StructDes, RP)
	return c
}

func (c *ChainTemplate) Handle() *ChainTemplate {
	if c.lock {
		panic("已生成的模板不可修改")
	}
	c.StructDes = append(c.StructDes, Handler)
	return c
}

func (c *ChainTemplate) Clone() *ChainTemplate {
	ret := NewChainTemplate()
	ret.StructDes = c.StructDes
	ret.Element = make([]any, len(c.Element), cap(c.Element))
	copy(ret.Element, c.Element)
	return ret
}

func (c *ChainTemplate) Fill(ele ...any) *ChainTemplate {
	c.Element = []any{}
	return c.Append(ele...)
}

func (c *ChainTemplate) Append(ele ...any) *ChainTemplate {
	c.lock = true
	for _, e := range ele {
		switch v := e.(type) {
		case string:
			c.Element = append(c.Element, v)
		case func(uint64, any, []any) (any, error):
			c.Element = append(c.Element, v)
		case *ChainTemplate:
			c.Element = append(c.Element, v.Element...)
		}
	}
	return c
}

func (c *ChainTemplate) Clear() *ChainTemplate {
	c.StructDes = []rune{}
	c.Element = []any{}
	c.nodeH = []*coordinator.Coordinator{}
	c.lock = false
	return c
}

func (c *ChainTemplate) Build() {
	var (
		index       = 0
		level       = 0
		currentNode = make(map[int]any)
		value       any
		ok          bool
		nr          any
	)

	helper := func(supply func() any) {
		value, ok = currentNode[level]
		if !ok {
			if len(c.nodeH) <= level {
				c.nodeH = append(c.nodeH, nil)
			}
			if c.nodeH[level] == nil {
				c.nodeH[level] = coordinator.NewCoordinator()
			}
			nr = c.nodeH[level].RootRouted
			if nr == nil {
				nr = supply()
				c.nodeH[level].RootRouted = nr
			}
		} else {
			switch r := value.(type) {
			case core.IRouter:
				switch vr := r.(type) {
				case *router.VersionRouter:
					nr, ok = vr.Has(c.Element[index].(string))
				default:
					nr, ok = r.Route(c.Element[index].(string))
				}
				if !ok {
					nr = supply()
					r.Add(c.Element[index].(string), nr)
				}
				index++
			case core.IHandler:
				nr, ok = r.Next()
				if !ok {
					nr = supply()
					h := r.(*handler.Handler)
					h.NextRouted = nr
				}
			case core.ICoordinator:
				nr, ok = r.Next()
				if !ok {
					nr = supply()
					co := r.(*coordinator.Coordinator)
					co.NextRouted = nr
				}
			}
		}
		currentNode[level] = nr
	}

	for _, des := range c.StructDes {
		switch des {
		case NameRouter:
			helper(func() any {
				return router.NewNameRouter()
			})
		case VersionRouter:
			helper(func() any {
				return router.NewVersionRouter()
			})
		case Handler:
			helper(func() any {
				return handler.NewHandler(c.Element[index+1].(func(uint64, any, []any) (any, error)))
			})
			index++
		case LP:
			helper(func() any {
				return coordinator.NewCoordinator()
			})
			if len(c.nodeH) <= level+1 {
				c.nodeH = append(c.nodeH, nil)
			}
			c.nodeH[level+1] = currentNode[level].(*coordinator.Coordinator)
			level++
		case RP:
			delete(currentNode, level)
			c.nodeH[level] = nil
			level--
		}
	}
	c.nodeH = []*coordinator.Coordinator{c.nodeH[0]}
}

func (c *ChainTemplate) Result() *ChainCaller {
	if len(c.nodeH) > 0 {
		return NewChainCaller(c.nodeH[0])
	}
	return nil
}

func (c *ChainTemplate) Description() []rune {
	return c.StructDes
}

func (c *ChainTemplate) Elem() []any {
	return c.Element
}

func HandlerHelper(h func(id uint64, result any, params []any) (any, error)) func(uint64, any, []any) (any, error) {
	return h
}
