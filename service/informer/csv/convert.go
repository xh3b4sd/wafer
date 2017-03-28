package csv

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	microerror "github.com/giantswarm/microkit/error"
	yaml "gopkg.in/yaml.v2"

	"github.com/xh3b4sd/wafer/service/informer"
	runtimeconfigdir "github.com/xh3b4sd/wafer/service/informer/csv/runtime/config/dir"
	runtimeconfigfile "github.com/xh3b4sd/wafer/service/informer/csv/runtime/config/file"
	runtimestatefile "github.com/xh3b4sd/wafer/service/informer/csv/runtime/state/file"
	statefileheader "github.com/xh3b4sd/wafer/service/informer/csv/runtime/state/file/header"
)

func fileToFiles(file runtimeconfigfile.File) ([]runtimestatefile.File, error) {
	files := []runtimestatefile.File{
		{
			Path: file.Path,
			Header: statefileheader.Header{
				Buy:    file.Header.Buy,
				Ignore: file.Header.Ignore,
				Sell:   file.Header.Sell,
				Time:   file.Header.Time,
			},
		},
	}

	return files, nil
}

func dirToFiles(dir runtimeconfigdir.Dir) ([]runtimestatefile.File, error) {
	var files []runtimestatefile.File

	outerFileInfos, err := ioutil.ReadDir(dir.Path)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}
	for _, ofi := range outerFileInfos {
		if !ofi.IsDir() {
			return nil, microerror.MaskAnyf(invalidExecutionError, "outer CSV dir must not contain files")
		}

		var stateFile runtimestatefile.File

		innerFileInfos, err := ioutil.ReadDir(filepath.Join(dir.Path, ofi.Name()))
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
		for _, ifi := range innerFileInfos {
			if ifi.IsDir() {
				return nil, microerror.MaskAnyf(invalidExecutionError, "outer CSV dir must not contain dirs")
			}

			if ifi.Name() == "chart.csv" {
				stateFile.Path = filepath.Join(dir.Path, ofi.Name(), ifi.Name())
				continue
			}

			if ifi.Name() == "header.yaml" {
				b, err := ioutil.ReadFile(filepath.Join(dir.Path, ofi.Name(), ifi.Name()))
				if err != nil {
					return nil, microerror.MaskAny(err)
				}

				var header statefileheader.Header
				err = yaml.Unmarshal(b, &header)
				if err != nil {
					return nil, microerror.MaskAny(err)
				}

				stateFile.Header.Buy = header.Buy
				stateFile.Header.Ignore = header.Ignore
				stateFile.Header.Sell = header.Sell
				stateFile.Header.Time = header.Time

				continue
			}

			return nil, microerror.MaskAnyf(invalidExecutionError, "additional file '%s' not allowed", ifi.Name())
		}

		files = append(files, stateFile)
	}

	return files, nil
}

func filesToPrices(files []runtimestatefile.File) ([][]informer.Price, error) {
	var prices [][]informer.Price

	for _, f := range files {
		// Read the CSV file.
		var fields [][]string
		{
			// TODO this should be abstracted with some more decent interface.
			csvFile, err := os.Open(f.Path)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}
			defer csvFile.Close()

			reader := csv.NewReader(csvFile)
			fields, err = reader.ReadAll()
			if err != nil {
				return nil, microerror.MaskAny(err)
			}
		}

		// Fill the prices channel.
		var ps []informer.Price
		{
			if f.Header.Ignore {
				fields = fields[1:]
			}

			for _, fs := range fields {
				b, err := strconv.ParseFloat(fs[f.Header.Buy], 64)
				if err != nil {
					return nil, microerror.MaskAny(err)
				}
				s, err := strconv.ParseFloat(fs[f.Header.Sell], 64)
				if err != nil {
					return nil, microerror.MaskAny(err)
				}
				t, err := strconv.ParseInt(fs[f.Header.Time], 10, 64)
				if err != nil {
					return nil, microerror.MaskAny(err)
				}

				price := informer.Price{
					Buy:  b,
					Sell: s,
					Time: time.Unix(t, 0),
				}

				ps = append(ps, price)
			}
		}

		prices = append(prices, ps)
	}

	return prices, nil
}
