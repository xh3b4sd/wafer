package corridor

import (
	microerror "github.com/giantswarm/microkit/error"
)

// Corridor describes the configuration of a price range in which buy events are
// allowed to happen. The values are provided in percent. Consider the following
// configuration.
//
//     Max     90
//
// This configuration means that buy events are not allowed to happen in case
// the inspected price is above 90% of the highest known price. The known price
// is taken from the observed chart window.
type Corridor struct {
	// Max is the maximum value within the allowed corridor.
	Max float64 `json:"max"`
}

func (c Corridor) Validate() error {
	if c.Max == 0 {
		return microerror.MaskAnyf(invalidConfigError, "Surge.Max must not be empty")
	}

	return nil
}
