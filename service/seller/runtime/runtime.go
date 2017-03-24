package runtime

import (
	"github.com/xh3b4sd/wafer/service/seller/runtime/config"
	"github.com/xh3b4sd/wafer/service/seller/runtime/state"
)

// Runtime holds information about the current runtime state of the seller used
// to judge stock market situations.
type Runtime struct {
	Config config.Config
	State  state.State
}
