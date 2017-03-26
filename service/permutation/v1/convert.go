package v1

import (
	"strconv"
	"strings"
	"time"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/spf13/cast"

	"github.com/xh3b4sd/wafer/service/permutation/runtime/config"
)

func IndizesToIndex(indizes []int) int {
	l := intsToStrings(indizes)
	s := strings.Join(l, "")
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return i
}

func IndexToIndizes(index int) []int {
	s := cast.ToString(index)
	l := strings.Split(s, "")
	i, err := stringsToInts(l)
	if err != nil {
		panic(err)
	}

	return i
}

func intsToStrings(ints []int) []string {
	var s []string

	for _, i := range ints {
		s = append(s, cast.ToString(i))
	}

	return s
}

func stringsToInts(strings []string) ([]int, error) {
	var i []int

	for _, s := range strings {
		n, err := cast.ToIntE(s)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
		i = append(i, n)
	}

	return i, nil
}

// maxFromConfigs creates a slice of numbers representing the maximum boundaries
// for shift permutations. Here the origin for the shift permutation is the list
// of permutation configs.
//
// NOTE maxFromConfigs converts several types to instances of type int. In case
// the original type exceeds the memory capacity of an int, maxFromConfigs
// returns the wrong result. Thus maxFromConfigs should only be used on maximum
// values that do not exceed the maximum memory capacity of int. This is
// 1705032704.
func maxFromConfigs(configs []config.Config) []int {
	var max []int

	for _, c := range configs {
		var m int

		switch c.Min.(type) {
		case time.Duration:
			d := (cast.ToDuration(c.Max) / cast.ToDuration(c.Step))
			if cast.ToDuration(c.Min) == 0 {
				m = int(d)
			} else {
				m = int(d) - 1
			}
		case float64:
			f := (cast.ToFloat64(c.Max) / cast.ToFloat64(c.Step))
			if cast.ToFloat64(c.Min) == 0 {
				m = int(f)
			} else {
				m = int(f) - 1
			}
		}

		max = append(max, m)
	}

	return max
}
