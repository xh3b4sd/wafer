package dir

import (
	microerror "github.com/giantswarm/microkit/error"
)

// Dir is the config for an absolute location of the CSV dir to consume. Inside
// this directory, at least one directory has to be resent. The name of this
// directory does not matter. The inner directories content matters. It has to
// contain two files. The chart data in CSV format, where the file name has to
// be chart.csv. The header options in YAML format, where the file name has to
// be header.yaml.
//
//    Dir
//    ├── 001
//    │   ├── chart.csv
//    │   └── header.yaml
//    ├── 002
//    │   ├── chart.csv
//    │   └── header.yaml
//    └── ...
//
type Dir struct {
	// Path is the absolute location of the CSV dir to consume.
	Path string
}

func (d Dir) Validate() error {
	if d.Path == "" {
		return microerror.MaskAnyf(invalidConfigError, "h.Path must not be empty")
	}

	return nil
}
