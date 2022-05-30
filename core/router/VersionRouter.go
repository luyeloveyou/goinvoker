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

func (vr *VersionRouter) Add(selector string, routed any) {
	vr.keys = append(vr.keys, selector)
	sort.Slice(vr.keys, func(i, j int) bool {
		if versionCompare(vr.keys[i], vr.keys[j]) < 0 {
			return true
		}
		return false
	})
	vr.cache[selector] = routed
}

func (vr *VersionRouter) Route(selector string) any {
	search := sort.Search(len(vr.keys), func(i int) bool {
		return versionCompare(vr.keys[i], selector) >= 0
	}) - 1
	if search < 0 || search > len(vr.keys) {
		return nil
	}
	return vr.cache[vr.keys[search]]
}