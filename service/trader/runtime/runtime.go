package runtime

import (
	"github.com/xh3b4sd/wafer/service/trader/runtime/config"
	"github.com/xh3b4sd/wafer/service/trader/runtime/state"
)

// Runtime holds information about the current runtime state of the trader.
type Runtime struct {
	Config config.Config
	State  state.State
}
