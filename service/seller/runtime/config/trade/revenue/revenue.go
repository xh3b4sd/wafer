package revenue

import (
	microerror "github.com/giantswarm/microkit/error"
)

type Revenue struct {
	// Min is the minimum revenue a single trade is allowed to make.
	Min float64 `json:"min"`
}

func (r Revenue) Validate() error {
	if r.Min == 0 {
		return microerror.MaskAnyf(invalidConfigError, "Revenue.Min must not be empty")
	}

	return nil
}
