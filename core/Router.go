package core

import (
	"math"
	"sort"
	"strings"
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

func (r *Router) Route(selector string) Object {
	return r.cache[selector]
}

func (r *Router) Add(selector string, routed Object) {
	r.cache[selector] = routed
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

func versionCompare(version1, version2 string) int {
	if len(version1) == 0 && len(version2) == 0 {
		return 0
	}
	if strings.EqualFold(version1, version2) {
		return 0
	}
	if len(version1) == 0 {
		return -1
	}
	if len(version2) == 0 {
		return 1
	}
	v1s := strings.Split(version1, ".")
	v2s := strings.Split(version2, ".")
	diff := 0
	minLength := int(math.Min(float64(len(v1s)), float64(len(v2s))))
	var (
		v1 string
		v2 string
	)
	for i := 0; i < minLength; i++ {
		v1 = v1s[i]
		v2 = v2s[i]
		diff = len(v1) - len(v2)
		if diff == 0 {
			diff = strings.Compare(v1, v2)
		}
		if diff != 0 {
			break
		}
	}
	if diff != 0 {
		return diff
	} else {
		return len(v1s) - len(v2s)
	}
}
