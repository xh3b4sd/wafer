package csv

import (
	"github.com/xh3b4sd/wafer/flag/service/informer/csv/index"
)

type CSV struct {
	File         string
	IgnoreHeader string
	Index        index.Index
}
