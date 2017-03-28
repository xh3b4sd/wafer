package file

import (
	"github.com/xh3b4sd/wafer/service/informer/csv/runtime/state/file/header"
)

type File struct {
	Header header.Header
	Path   string
}
