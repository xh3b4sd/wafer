package state

import (
	"github.com/xh3b4sd/wafer/service/analyzer/runtime/state/config"
	"github.com/xh3b4sd/wafer/service/analyzer/runtime/state/informer"
	"github.com/xh3b4sd/wafer/service/analyzer/runtime/state/permutation"
)

type State struct {
	Config      config.Config           `json:"config"`
	Informer    informer.Informer       `json:"informer"`
	Permutation permutation.Permutation `json:"permutation"`
}
