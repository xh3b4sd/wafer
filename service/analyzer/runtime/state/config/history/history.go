package history

import (
	"github.com/xh3b4sd/wafer/service/analyzer/runtime/config"
)

type History struct {
	Config   config.Config `json:"config"`
	Cycles   []int64       `json:"cycles"`
	Indizes  []int         `json:"indizes"`
	Revenues []float64     `json:"revenues"`
}
