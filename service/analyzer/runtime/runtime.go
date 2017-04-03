package runtime

import (
	"github.com/xh3b4sd/wafer/service/analyzer/runtime/config"
	"github.com/xh3b4sd/wafer/service/analyzer/runtime/state"
)

// Runtime holds information about the current runtime of the analyzer used to
// judge stock market situations.
type Runtime struct {
	Config *config.Config
	State  state.State
}
