package flag

import (
	"github.com/giantswarm/microkit/flag"

	"github.com/xh3b4sd/wafer/flag/service"
)

type Flag struct {
	Service service.Service
}

func New() *Flag {
	f := &Flag{}
	flag.Init(f)
	return f
}
