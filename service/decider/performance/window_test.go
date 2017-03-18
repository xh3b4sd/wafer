package performance

import (
	"reflect"
	"testing"
	"time"

	"github.com/xh3b4sd/wafer/service/informer"
)

func Test_calculateWindow(t *testing.T) {
	testCases := []struct {
		HistWindow   []informer.Price
		ConfigWindow time.Duration
		Expected     []informer.Price
	}{
		{
			HistWindow: []informer.Price{
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(1, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(2, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(3, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(4, 0)},
			},
			ConfigWindow: 3 * time.Second,
			Expected: []informer.Price{
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(2, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(3, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(4, 0)},
			},
		},
		{
			HistWindow: []informer.Price{
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(4, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(5, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(6, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(7, 0)},
			},
			ConfigWindow: 2 * time.Second,
			Expected: []informer.Price{
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(6, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(7, 0)},
			},
		},
		{
			HistWindow: []informer.Price{
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(6, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(8, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(10, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(12, 0)},
			},
			ConfigWindow: 5 * time.Second,
			Expected: []informer.Price{
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(8, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(10, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(12, 0)},
			},
		},
	}

	for i, testCase := range testCases {
		expected, err := calculateWindow(testCase.HistWindow, testCase.ConfigWindow)
		if err != nil {
			t.Fatal("case", i+1, "expected", nil, "got", err)
		}
		if !reflect.DeepEqual(expected, testCase.Expected) {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", expected)
		}
	}
}
