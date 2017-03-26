package permutation

import (
	"github.com/xh3b4sd/wafer/service/permutation/runtime/config"
)

type Object interface {
	GetPermConfigs() []config.Config
	SetPermValue(permID string, permValue interface{}) error
}

type Permutation interface {
	ValueFor(indizes []int) (interface{}, error)
}
