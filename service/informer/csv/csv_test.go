package csv

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/xh3b4sd/wafer/service/informer"
)

func Test_New(t *testing.T) {
	testCases := []struct {
		BuyIndex     int
		FileFunc     func() string
		IgnoreHeader bool
		Expected     map[int]informer.Price
		SellIndex    int
		TimeIndex    int
	}{
		{
			BuyIndex: 9,
			FileFunc: func() string {
				f, err := filepath.Abs("./fixtures/001.csv")
				if err != nil {
					t.Fatal("expected", nil, "got", err)
				}
				return f
			},
			IgnoreHeader: true,
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
			SellIndex: 10,
			TimeIndex: 12,
		},
	}

	for i, testCase := range testCases {
		newConfig := DefaultConfig()
		newConfig.BuyIndex = testCase.BuyIndex
		newConfig.File = testCase.FileFunc()
		newConfig.IgnoreHeader = testCase.IgnoreHeader
		newConfig.SellIndex = testCase.SellIndex
		newConfig.TimeIndex = testCase.TimeIndex

		newInformer, err := New(newConfig)
		if err != nil {
			t.Fatal("case", i+1, "expected", nil, "got", err)
		}

		var j int
		for price := range newInformer.Prices() {
			for k, p := range testCase.Expected {
				if j != k {
					continue
				}

				if !reflect.DeepEqual(price, p) {
					fmt.Printf("%#v\n", price)
					t.Fatal("case", i+1, "expected", p, "got", price)
				}
			}

			j++
		}
	}
}
