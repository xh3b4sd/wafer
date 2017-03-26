package informer

import (
	"github.com/xh3b4sd/wafer/flag/service/informer/csv"
)

type Informer struct {
	CSV  csv.CSV
	Kind string
}
