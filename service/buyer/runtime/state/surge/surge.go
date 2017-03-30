package surge

import (
	"github.com/xh3b4sd/wafer/service/informer"
)

type Surge struct {
	// Last is a list of consecutive price events representing the last surge
	// within the chart window.
	Last []informer.Price
}
