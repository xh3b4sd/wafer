package informer

import (
	"github.com/xh3b4sd/wafer/service/informer/csv/runtime/state/price"
)

type Informer struct {
	Prices []price.Price `json:"prices"`
}
