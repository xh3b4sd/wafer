package flag

import (
	microerror "github.com/giantswarm/microkit/error"

	"github.com/xh3b4sd/wafer/command/render/flag/index"
	"github.com/xh3b4sd/wafer/command/render/flag/server"
)

type Flag struct {
	File         string
	IgnoreHeader bool
	Index        index.Index
	Server       server.Server
}

func (f *Flag) Validate() error {
	if f.File == "" {
		return microerror.MaskAnyf(invalidFlagsError, "flag.File must not be empty")
	}
	if f.Server.Port == 0 {
		return microerror.MaskAnyf(invalidFlagsError, "flag.Server.Port must not be empty")
	}

	return nil
}
