package v1

import (
	"fmt"
)

func typesEqual(types []interface{}) bool {
	if len(types) == 0 {
		return false
	}

	for _, t := range types {
		if fmt.Sprintf("%T", t) != fmt.Sprintf("%T", types[0]) {
			return false
		}
	}

	return true
}

func indizesGreaterThan(indizes, max []int) bool {
	if len(indizes) > len(max) {
		return true
	}

	for i, _ := range indizes {
		if indizes[i] > max[i] {
			return true
		}
	}

	return false
}
