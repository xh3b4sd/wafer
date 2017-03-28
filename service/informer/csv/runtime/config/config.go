package config

import (
	"github.com/xh3b4sd/wafer/service/informer/csv/runtime/config/dir"
	"github.com/xh3b4sd/wafer/service/informer/csv/runtime/config/file"
)

type Config struct {
	Dir  dir.Dir
	File file.File
}
