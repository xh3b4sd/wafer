package csv

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/xh3b4sd/wafer/service/informer"
	runtimeconfigdir "github.com/xh3b4sd/wafer/service/informer/csv/runtime/config/dir"
	runtimeconfigfile "github.com/xh3b4sd/wafer/service/informer/csv/runtime/config/file"
	configfileheader "github.com/xh3b4sd/wafer/service/informer/csv/runtime/config/file/header"
)

func Test_Informer_File_Prices(t *testing.T) {
	path, err := filepath.Abs("./fixtures/file/001.csv")
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	testCases := []struct {
		File     runtimeconfigfile.File
		Expected map[int]informer.Price
	}{
		{
			File: runtimeconfigfile.File{
				Header: configfileheader.Header{
					Buy:    9,
					Ignore: true,
					Sell:   10,
					Time:   12,
				},
				Path: path,
			},
			Expected: map[int]informer.Price{
				0: {
					Buy:  797.4000000000,
					Sell: 797.0000000000,
					Time: time.Unix(1391212802, 0),
				},
				4: {
					Buy:  798.8900000000,
					Sell: 797.0000000000,
					Time: time.Unix(1391213041, 0),
				},
				9: {
					Buy:  796.9000000000,
					Sell: 793.0000000000,
					Time: time.Unix(1391213342, 0),
				},
			},
		},
	}

	for i, testCase := range testCases {
		newConfig := DefaultConfig()
		newConfig.File = testCase.File
		newInformer, err := New(newConfig)
		if err != nil {
			t.Fatal("case", i+1, "expected", nil, "got", err)
		}

		// The file configuration for the CSV informer must result in one channel in
		// the channel list, because there is only one CSV file to parse.
		prices := newInformer.Prices()
		if len(prices) != 1 {
			t.Fatal("case", i+1, "expected", 1, "got", len(prices))
		}

		var j int
		for price := range prices[0] {
			for k, p := range testCase.Expected {
				if j != k {
					continue
				}

				if !reflect.DeepEqual(price, p) {
					t.Fatal("case", i+1, "expected", p, "got", price)
				}
			}

			j++
		}
	}
}

// Test_Informer_File_Prices_MultipleCalls makes sure Informer.Prices can be
// called mutliple times, which means the actual price data is cached and always
// available in memory.
func Test_Informer_File_Prices_MultipleCalls(t *testing.T) {
	path, err := filepath.Abs("./fixtures/file/001.csv")
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	newConfig := DefaultConfig()
	newConfig.File.Header.Buy = 9
	newConfig.File.Header.Ignore = true
	newConfig.File.Header.Sell = 10
	newConfig.File.Header.Time = 12
	newConfig.File.Path = path

	newInformer, err := New(newConfig)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// The file configuration for the CSV informer must result in one channel in
	// the channel list, because there is only one CSV file to parse.
	prices1 := newInformer.Prices()
	if len(prices1) != 1 {
		t.Fatal("expected", 1, "got", len(prices1))
	}
	p1 := <-prices1[0]
	prices2 := newInformer.Prices()
	if len(prices2) != 1 {
		t.Fatal("expected", 1, "got", len(prices2))
	}
	p2 := <-prices2[0]
	prices3 := newInformer.Prices()
	if len(prices3) != 1 {
		t.Fatal("expected", 1, "got", len(prices3))
	}
	p3 := <-prices3[0]

	if !reflect.DeepEqual(p1, p2) {
		t.Fatal("expected", true, "got", false)
	}
	if !reflect.DeepEqual(p1, p3) {
		t.Fatal("expected", true, "got", false)
	}
	if !reflect.DeepEqual(p2, p3) {
		t.Fatal("expected", true, "got", false)
	}
}

func Test_Informer_Dir_Prices(t *testing.T) {
	path, err := filepath.Abs("./fixtures/dir/")
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	testCases := []struct {
		Dir      runtimeconfigdir.Dir
		Expected []map[int]informer.Price
	}{
		{
			Dir: runtimeconfigdir.Dir{
				Path: path,
			},
			Expected: []map[int]informer.Price{
				{
					0: {
						Buy:  797.4000000000,
						Sell: 797.0000000000,
						Time: time.Unix(1391212802, 0),
					},
					4: {
						Buy:  798.8900000000,
						Sell: 797.0000000000,
						Time: time.Unix(1391213041, 0),
					},
					9: {
						Buy:  796.9000000000,
						Sell: 793.0000000000,
						Time: time.Unix(1391213342, 0),
					},
				},
				{
					0: {
						Buy:  797.4000000000,
						Sell: 797.0000000000,
						Time: time.Unix(1391212802, 0),
					},
					4: {
						Buy:  798.8900000000,
						Sell: 797.0000000000,
						Time: time.Unix(1391213041, 0),
					},
					9: {
						Buy:  796.9000000000,
						Sell: 793.0000000000,
						Time: time.Unix(1391213342, 0),
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		newConfig := DefaultConfig()
		newConfig.Dir = testCase.Dir
		newInformer, err := New(newConfig)
		if err != nil {
			t.Fatal("case", i+1, "expected", nil, "got", err)
		}

		// The dir configuration for the CSV informer must result in two channels in
		// the channel list, because there are two CSV directories to parse.
		prices := newInformer.Prices()
		if len(prices) != 2 {
			t.Fatal("case", i+1, "expected", 2, "got", len(prices))
		}

		for d, c := range prices {
			var j int
			for price := range c {
				for k, p := range testCase.Expected[d] {
					if j != k {
						continue
					}

					if !reflect.DeepEqual(price, p) {
						t.Fatal("case", i+1, "expected", p, "got", price)
					}
				}

				j++
			}
		}
	}
}
