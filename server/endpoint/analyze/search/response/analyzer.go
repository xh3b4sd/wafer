package response

import (
	"github.com/xh3b4sd/wafer/service/analyzer/runtime/state"
)

type Analyzer struct {
	State state.State `json:"state"`
}
