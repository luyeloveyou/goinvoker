package core

import (
	"math"
	"strings"
)

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
