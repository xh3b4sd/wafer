package v1

import (
	"time"

	microerror "github.com/giantswarm/microkit/error"
)

// DurationConfig is the configuration used to create a new duration.
type DurationConfig struct {
	// Settings.
	Min  time.Duration
	Max  time.Duration
	Step time.Duration
}

// DefaultDurationConfig returns the default configuration used to create a new
// duration by best effort.
func DefaultDurationConfig() DurationConfig {
	return DurationConfig{
		// Settings.
		Min:  0,
		Max:  0,
		Step: 0,
	}
}

// NewDuration creates a new configured duration.
func NewDuration(config DurationConfig) (*Duration, error) {
	// Settings.
	if config.Max == 0 {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Max must not be empty")
	}
	if config.Step == 0 {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Step must not be empty")
	}

	newDuration := &Duration{
		// Settings.
		min:  config.Min,
		max:  config.Max,
		step: config.Step,
	}

	return newDuration, nil
}

// Duration implements permutation.Permutation.
type Duration struct {
	// Settings.
	min  time.Duration
	max  time.Duration
	step time.Duration
}

func (d *Duration) ValueFor(indizes []int) (interface{}, error) {
	if len(indizes) != 1 {
		return 0, microerror.MaskAnyf(invalidExecutionError, "indizes must have length 1")
	}
	index := indizes[0]

	var value time.Duration

	if index == 0 {
		value = d.min
	} else {
		value = d.min + (time.Duration(index) * d.step)
	}

	if value > d.max {
		return 0, microerror.MaskAny(invalidExecutionError)
	}

	return value, nil
}
