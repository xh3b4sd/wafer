package buy

import (
	"github.com/xh3b4sd/wafer/service/informer"
)

type Buy struct {
	// Last is the price event of the last buy that took place.
	Last informer.Price
}
