package runtime

import (
	"github.com/xh3b4sd/wafer/service/permutation/runtime/config"
)

// Runtime holds information about the current runtime state of the permutation.
type Runtime struct {
	Configs []config.Config
}
