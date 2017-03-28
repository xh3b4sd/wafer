package runtime

import (
	"github.com/xh3b4sd/wafer/service/informer/csv/runtime/config"
	"github.com/xh3b4sd/wafer/service/informer/csv/runtime/state"
)

type Runtime struct {
	Config config.Config
	State  state.State
}
