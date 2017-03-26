package v1

import (
	microerror "github.com/giantswarm/microkit/error"
)

// Float64Config is the configuration used to create a new float64.
type Float64Config struct {
	// Settings.
	Min  float64
	Max  float64
	Step float64
}

// DefaultFloat64Config returns the default configuration used to create a new
// float64 by best effort.
func DefaultFloat64Config() Float64Config {
	return Float64Config{
		// Settings.
		Min:  0,
		Max:  0,
		Step: 0,
	}
}

// NewFloat64 creates a new configured float64.
func NewFloat64(config Float64Config) (*Float64, error) {
	// Settings.
	if config.Max == 0 {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Max must not be empty")
	}
	if config.Step == 0 {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Step must not be empty")
	}

	newFloat64 := &Float64{
		// Settings.
		min:  config.Min,
		max:  config.Max,
		step: config.Step,
	}

	return newFloat64, nil
}

// Float64 implements permutation.Permutation.
type Float64 struct {
	// Settings.
	min  float64
	max  float64
	step float64
}

func (f *Float64) ValueFor(indizes []int) (interface{}, error) {
	if len(indizes) != 1 {
		return 0, microerror.MaskAnyf(invalidExecutionError, "indizes must have length 1")
	}
	index := indizes[0]

	var value float64

	if index == 0 {
		value = f.min
	} else {
		value = f.min + (float64(index) * f.step)
	}

	if value > f.max {
		return 0, microerror.MaskAny(invalidExecutionError)
	}

	return value, nil
}
