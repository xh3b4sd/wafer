package history

import (
	"github.com/xh3b4sd/wafer/service/analyzer/runtime/config"
)

type History struct {
	Config  config.Config
	Indizes []int
	Revenue float64
}
