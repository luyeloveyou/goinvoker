package router

import "sort"

type VersionRouter struct {
	cache map[string]any
	keys  []string
}

func NewVersionRouter() *VersionRouter {
	return &VersionRouter{
		cache: make(map[string]any),
		keys:  []string{},
	}
}

func (vr *VersionRouter) Add(selector string, routed any) bool {
	vr.keys = append(vr.keys, selector)
	sort.Slice(vr.keys, func(i, j int) bool {
		if versionCompare(vr.keys[i], vr.keys[j]) >= 0 {
			return true
		}
		return false
	})
	vr.cache[selector] = routed
	return true
}

func (vr *VersionRouter) Route(selector string) (any, bool) {
	search := sort.Search(len(vr.keys), func(i int) bool {
		return versionCompare(vr.keys[i], selector) <= 0
	})
	if search >= len(vr.keys) {
		return nil, false
	}
	v, ok := vr.cache[vr.keys[search]]
	return v, ok
}
