package csv

import (
	"github.com/xh3b4sd/wafer/flag/service/informer/csv/header"
)

type CSV struct {
	Dir    string
	File   string
	Header header.Header
}
