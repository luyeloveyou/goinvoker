package core

import (
	"sort"
)

type IRouter interface {
	Route(selector string) Object
	Add(selector string, routed Object)
}

type Router struct {
	cache map[string]Object
}

func NewRouter() *Router {
	return &Router{cache: make(map[string]Object)}
}

type VersionRouter struct {
	cache map[string]Object
	keys  []string
}

func NewVersionRouter() *VersionRouter {
	return &VersionRouter{
		cache: make(map[string]Object),
		keys:  []string{},
	}
}

func (r *Router) Route(selector string) Object {
	return r.cache[selector]
}

func (r *Router) Add(selector string, routed Object) {
	r.cache[selector] = routed
}

func (vr *VersionRouter) Add(selector string, routed Object) {
	vr.keys = append(vr.keys, selector)
	sort.Slice(vr.keys, func(i, j int) bool {
		if versionCompare(vr.keys[i], vr.keys[j]) < 0 {
			return true
		}
		return false
	})
	vr.cache[selector] = routed
}

func (vr *VersionRouter) Route(selector string) Object {
	search := sort.Search(len(vr.keys), func(i int) bool {
		return versionCompare(vr.keys[i], selector) >= 0
	}) - 1
	if search < 0 || search > len(vr.keys) {
		return nil
	}
	return vr.cache[vr.keys[search]]
}
