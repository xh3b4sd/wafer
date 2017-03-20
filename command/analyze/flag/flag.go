package flag

import (
	microerror "github.com/giantswarm/microkit/error"

	"github.com/xh3b4sd/wafer/command/analyze/flag/index"
)

type Flag struct {
	File         string
	IgnoreHeader bool
	Index        index.Index
}

func (f *Flag) Validate() error {
	if f.File == "" {
		return microerror.MaskAnyf(invalidFlagsError, "flag.File must not be empty")
	}

	return nil
}
