package v1

import (
	microerror "github.com/giantswarm/microkit/error"
)

func ShiftIndizes(indizes, max []int) ([]int, error) {
	if len(indizes) == 0 {
		return nil, microerror.MaskAnyf(invalidExecutionError, "indizes must not be empty")
	}
	if len(max) == 0 {
		return nil, microerror.MaskAnyf(invalidExecutionError, "max must not be empty")
	}
	if len(indizes) != len(max) {
		return nil, microerror.MaskAnyf(invalidExecutionError, "indizes must be as long as max")
	}

	var copiedIndizes []int
	for _, i := range indizes {
		copiedIndizes = append(copiedIndizes, i)
	}

	copiedIndizes[0] += 1

	for i, _ := range copiedIndizes {
		if copiedIndizes[i] > max[i] {
			copiedIndizes[i] = 0

			if i+1 < len(copiedIndizes) {
				copiedIndizes[i+1] += 1
			}
		}
	}

	return copiedIndizes, nil
}
