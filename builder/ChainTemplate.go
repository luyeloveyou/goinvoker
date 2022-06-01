package builder

import (
	"github.com/luyeloveyou/goinvoker/core"
	"github.com/luyeloveyou/goinvoker/core/coordinator"
	"github.com/luyeloveyou/goinvoker/core/handler"
	"github.com/luyeloveyou/goinvoker/core/router"
	"strings"
)

const (
	NameRouter    = 'N'
	VersionRouter = 'V'
	Handler       = 'H'
	LP            = '('
	RP            = ')'
)

type IChainTemplate interface {
	Name() IChainTemplate
	Version() IChainTemplate
	Next(template IChainTemplate) IChainTemplate
	Handle() IChainTemplate
	Fill(ele ...any) IChainTemplate
	Append(ele ...any) IChainTemplate
	Build()
	Coordinator() core.ICoordinator
	Clone() IChainTemplate
	Clear() IChainTemplate
	Description() []rune
	Path() string
}

type ChainTemplate struct {
	StructDes []rune
	Paths     string
	Handlers  []func(uint64, any, []any) (any, error)
	nodeH     []*coordinator.Coordinator
	lock      bool
}

func NewChainTemplate() *ChainTemplate {
	return &ChainTemplate{
		StructDes: []rune{},
		Paths:     "",
		Handlers:  make([]func(uint64, any, []any) (any, error), 0, 4),
		nodeH:     make([]*coordinator.Coordinator, 0, 4),
		lock:      false,
	}
}

func (c *ChainTemplate) Name() IChainTemplate {
	if c.lock {
		panic("已生成的模板不可修改")
	}
	c.StructDes = append(c.StructDes, NameRouter)
	return c
}

func (c *ChainTemplate) Version() IChainTemplate {
	if c.lock {
		panic("已生成的模板不可修改")
	}
	c.StructDes = append(c.StructDes, VersionRouter)
	return c
}

func (c *ChainTemplate) Next(template IChainTemplate) IChainTemplate {
	if c.lock {
		panic("已生成的模板不可修改")
	}
	description := template.Description()
	c.StructDes = append(c.StructDes, LP)
	c.StructDes = append(c.StructDes, description...)
	c.StructDes = append(c.StructDes, RP)
	return c
}

func (c *ChainTemplate) Handle() IChainTemplate {
	if c.lock {
		panic("已生成的模板不可修改")
	}
	c.StructDes = append(c.StructDes, Handler)
	return c
}

func (c *ChainTemplate) Clone() IChainTemplate {
	ret := NewChainTemplate()
	ret.StructDes = c.StructDes
	ret.Paths = c.Paths
	for _, handler := range c.Handlers {
		ret.Handlers = append(ret.Handlers, handler)
	}
	return ret
}

func (c *ChainTemplate) Fill(ele ...any) IChainTemplate {
	c.lock = true
	path := make([]string, 0, 4)
	c.Handlers = append([]func(uint64, any, []any) (any, error){})
	for _, e := range ele {
		switch v := e.(type) {
		case string:
			path = append(path, v)
		case func(uint64, any, []any) (any, error):
			c.Handlers = append(c.Handlers, v)
		case IChainTemplate:
			t := v.(*ChainTemplate)
			path = append(path, strings.Split(t.Paths, "|")...)
			c.Handlers = append(c.Handlers, t.Handlers...)
		}
	}
	c.Paths = strings.Join(path, "|")
	return c
}

func (c *ChainTemplate) Append(ele ...any) IChainTemplate {
	c.lock = true
	path := []string{c.Paths}
	for _, e := range ele {
		switch v := e.(type) {
		case string:
			path = append(path, v)
		case func(uint64, any, []any) (any, error):
			c.Handlers = append(c.Handlers, v)
		case IChainTemplate:
			t := v.(*ChainTemplate)
			path = append(path, strings.Split(t.Paths, "|")...)
			c.Handlers = append(c.Handlers, t.Handlers...)
		}
	}
	c.Paths = strings.Join(path, "|")
	return c
}

func (c *ChainTemplate) Clear() IChainTemplate {
	c.StructDes = []rune{}
	c.Paths = ""
	c.Handlers = []func(uint64, any, []any) (any, error){}
	c.nodeH = []*coordinator.Coordinator{}
	c.lock = false
	return c
}

func (c *ChainTemplate) Build() {
	var (
		path         = strings.Split(c.Paths, "|")
		pathIndex    = 0
		handlerIndex = 0
		level        = 0
		currentNode  = make(map[int]any)
		value        any
		ok           bool
		nr           any
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
					nr, ok = vr.Has(path[pathIndex])
				default:
					nr, ok = r.Route(path[pathIndex])
				}
				if !ok {
					nr = supply()
					r.Add(path[pathIndex], nr)
				}
				pathIndex++
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
				handlerIndex++
				return handler.NewHandler(c.Handlers[handlerIndex-1])
			})
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
}

func (c *ChainTemplate) Coordinator() core.ICoordinator {
	if len(c.nodeH) > 0 {
		return c.nodeH[0]
	}
	return nil
}

func (c *ChainTemplate) Description() []rune {
	return c.StructDes
}

func (c *ChainTemplate) Path() string {
	return c.Paths
}

func HandlerHelper(h func(id uint64, result any, params []any) (any, error)) func(uint64, any, []any) (any, error) {
	return h
}
