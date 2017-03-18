package performance

import (
	"github.com/juju/errgo"
)

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var notEnoughDataError = errgo.New("not enough data")

// IsNotEnoughData asserts notEnoughDataError.
func IsNotEnoughData(err error) bool {
	return errgo.Cause(err) == notEnoughDataError
}
