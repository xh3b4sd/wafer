package history

import (
	"github.com/xh3b4sd/wafer/service/analyzer/runtime/config"
)

type History struct {
	Config  config.Config `json:"config"`
	Indizes []int         `json:"indizes"`
	Revenue float64       `json:"revenue"`
}
