package informer

import (
	"github.com/xh3b4sd/wafer/service/analyzer/runtime/state/informer/price"
)

type Informer struct {
	Prices []price.Price `json:"price"`
}
