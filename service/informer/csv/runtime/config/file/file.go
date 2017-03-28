package file

import (
	microerror "github.com/giantswarm/microkit/error"

	"github.com/xh3b4sd/wafer/service/informer/csv/runtime/config/file/header"
)

type File struct {
	Header header.Header
	Path   string
}

func (f File) Validate() error {
	if f.Path == "" {
		return microerror.MaskAnyf(invalidConfigError, "f.Path must not be empty")
	}

	err := f.Header.Validate()
	if err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}
