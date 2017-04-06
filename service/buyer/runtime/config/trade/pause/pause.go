package pause

import (
	"time"

	microerror "github.com/giantswarm/microkit/error"
)

type Pause struct {
	// Min is the minimum time to wait between buys.
	Min time.Duration `json:"min"`
}

func (p Pause) Validate() error {
	if p.Min.Seconds() == 0 {
		return microerror.MaskAnyf(invalidConfigError, "Surge.Min must not be empty")
	}

	return nil
}
