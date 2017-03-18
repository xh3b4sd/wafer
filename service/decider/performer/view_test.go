package performer

import (
	"reflect"
	"testing"
	"time"

	"github.com/xh3b4sd/wafer/service/informer"
)

func Test_calculateViews(t *testing.T) {
	testCases := []struct {
		HistWindow        []informer.Price
		ConfigView        time.Duration
		ConfigConvolution time.Duration
		Expected          []View
	}{
		{
			HistWindow: []informer.Price{
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(1, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(2, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(3, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(4, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(5, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(6, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(7, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(8, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(9, 0)},
			},
			ConfigView:        5 * time.Second,
			ConfigConvolution: 2 * time.Second,
			Expected: []View{
				{
					HasLeftNeighbour:  false,
					HasRightNeighbour: true,
					LeftBound:         informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(1, 0)},
					RightBound:        informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(5, 0)},
				},
				{
					HasLeftNeighbour:  true,
					HasRightNeighbour: true,
					LeftBound:         informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(3, 0)},
					RightBound:        informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(7, 0)},
				},
				{
					HasLeftNeighbour:  true,
					HasRightNeighbour: false,
					LeftBound:         informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(5, 0)},
					RightBound:        informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(9, 0)},
				},
			},
		},
		{
			HistWindow: []informer.Price{
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(1, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(2, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(3, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(4, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(5, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(6, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(7, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(8, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(9, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(10, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(11, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(12, 0)},
			},
			ConfigView:        5 * time.Second,
			ConfigConvolution: 3 * time.Second,
			Expected: []View{
				{
					HasLeftNeighbour:  false,
					HasRightNeighbour: true,
					LeftBound:         informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(1, 0)},
					RightBound:        informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(5, 0)},
				},
				{
					HasLeftNeighbour:  true,
					HasRightNeighbour: true,
					LeftBound:         informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(4, 0)},
					RightBound:        informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(8, 0)},
				},
				{
					HasLeftNeighbour:  true,
					HasRightNeighbour: false,
					LeftBound:         informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(7, 0)},
					RightBound:        informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(11, 0)},
				},
			},
		},
		{
			HistWindow: []informer.Price{
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(1, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(3, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(3, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(5, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(5, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(6, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(11, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(13, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(15, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(27, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(29, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(30, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(31, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(32, 0)},
				{Buy: 23.45, Sell: 23.45, Time: time.Unix(34, 0)},
			},
			ConfigView:        5 * time.Second,
			ConfigConvolution: 3 * time.Second,
			Expected: []View{
				{
					HasLeftNeighbour:  false,
					HasRightNeighbour: false,
					LeftBound:         informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(1, 0)},
					RightBound:        informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(5, 0)},
				},
				{
					HasLeftNeighbour:  false,
					HasRightNeighbour: false,
					LeftBound:         informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(11, 0)},
					RightBound:        informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(15, 0)},
				},
				{
					HasLeftNeighbour:  false,
					HasRightNeighbour: true,
					LeftBound:         informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(27, 0)},
					RightBound:        informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(31, 0)},
				},
				{
					HasLeftNeighbour:  true,
					HasRightNeighbour: false,
					LeftBound:         informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(30, 0)},
					RightBound:        informer.Price{Buy: 23.45, Sell: 23.45, Time: time.Unix(34, 0)},
				},
			},
		},
	}

	for i, testCase := range testCases {
		expected, err := calculateViews(testCase.HistWindow, testCase.ConfigView, testCase.ConfigConvolution)
		if err != nil {
			t.Fatal("case", i+1, "expected", nil, "got", err)
		}
		if !reflect.DeepEqual(expected, testCase.Expected) {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", expected)
		}
	}
}
