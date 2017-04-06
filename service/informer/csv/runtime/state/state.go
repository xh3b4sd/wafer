package state

import (
	"github.com/xh3b4sd/wafer/service/informer/csv/runtime/state/file"
	"github.com/xh3b4sd/wafer/service/informer/csv/runtime/state/price"
)

type State struct {
	Files  []file.File   `json:"files"`
	Prices []price.Price `json:"price"`
}
