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
	Infinite      = '*'
	LP            = '('
	RP            = ')'
)

type ChainTemplate struct {
	StructDes []rune
	Element   []any
	result    *coordinator.Coordinator
	lock      bool
}

func NewChainTemplate() *ChainTemplate {
	return &ChainTemplate{
		StructDes: []rune{},
		Element:   []any{},
		result:    nil,
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

func (c *ChainTemplate) Next(template *ChainTemplate, n int) *ChainTemplate {
	if c.lock {
		panic("已生成的模板不可修改")
	}
	if n > 0 {
		for i := 0; i < n; i++ {
			description := template.Description()
			c.StructDes = append(c.StructDes, LP)
			c.StructDes = append(c.StructDes, description...)
			c.StructDes = append(c.StructDes, RP)
		}
	} else {
		description := template.Description()
		c.StructDes = append(c.StructDes, LP)
		c.StructDes = append(c.StructDes, description...)
		c.StructDes = append(c.StructDes, RP, Infinite)
	}
	c.lock = true
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
	ret.StructDes = make([]rune, len(c.StructDes), cap(c.StructDes))
	ret.Element = make([]any, len(c.Element), cap(c.Element))
	copy(ret.StructDes, c.StructDes)
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
	c.result = nil
	c.lock = false
	return c
}

type _stack struct {
	node    *coordinator.Coordinator
	current any
}

func (c *ChainTemplate) Build() {
	var (
		desIndex   = 0
		eleIndex   = 0
		stackIndex = 0
		ok         bool
		nr         any
	)

	stack := make(map[int]*_stack)

	newStack := func(init *coordinator.Coordinator) {
		if init == nil {
			init = coordinator.NewCoordinator()
		}
		stack[stackIndex] = &_stack{
			node:    init,
			current: init,
		}
	}

	newStack(c.result)

	helper := func(supply func() any) {
		s := stack[stackIndex]
		if s.node == s.current {
			nr = s.node.RootRouted
			if nr == nil {
				nr = supply()
				s.node.RootRouted = nr
			}
		} else {
			switch r := s.current.(type) {
			case core.IRouter:
				selector := c.Element[eleIndex].(string)
				nr, ok = r.Has(selector)
				if !ok {
					nr = supply()
					r.Add(selector, nr)
				}
				eleIndex++
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
		s.current = nr
	}

	for eleIndex < len(c.Element) {
		des := c.StructDes[desIndex]
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
				return handler.NewHandler(c.Element[eleIndex+1].(func(uint64, any, []any) (any, error)))
			})
			eleIndex++
		case LP:
			helper(func() any {
				return coordinator.NewCoordinator()
			})
			stackIndex++
			newStack(stack[stackIndex-1].current.(*coordinator.Coordinator))
		case RP:
			delete(stack, stackIndex)
			stackIndex--
			if desIndex+1 < len(c.StructDes) && eleIndex < len(c.Element) && c.StructDes[desIndex+1] == Infinite {
				for ; c.StructDes[desIndex] != LP; desIndex-- {
				}
				desIndex--
			}
		}
		desIndex++
	}
	c.result = stack[0].node
	for i := 0; i < stackIndex; i++ {
		delete(stack, stackIndex)
	}
}

func (c *ChainTemplate) Result() *ChainCaller {
	if c.result != nil {
		return NewChainCaller(c.result)
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
