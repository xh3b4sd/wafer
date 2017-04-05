package surge

import (
	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/surge/duration"
)

type Surge struct {
	Duration duration.Duration `json:"duration"`
	// Min is the minimum angle the observed chart is allowed to have before a buy
	// event can happen.
	Min float64 `json:"min"`
	// Tolerance is the percentage allowed to vary when calculating a chart's
	// surge.
	Tolerance float64 `json:"tolerance"`
}
