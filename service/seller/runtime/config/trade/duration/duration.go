package duration

import (
	"time"

	microerror "github.com/giantswarm/microkit/error"
)

type Duration struct {
	// Min is the minimum time a single trade is allowed to take.
	Min time.Duration `json:"min"`
}

func (d Duration) Validate() error {
	if d.Min.Seconds() == 0 {
		return microerror.MaskAnyf(invalidConfigError, "Duration.Min must not be empty")
	}

	return nil
}
