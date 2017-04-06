package v1

import (
	"testing"
	"time"

	"github.com/xh3b4sd/wafer/service/buyer/runtime"
	"github.com/xh3b4sd/wafer/service/informer"
)

func Test_IsInsideMinTradePause(t *testing.T) {
	testCases := []struct {
		RuntimeFunc func() runtime.Runtime
		Expected    bool
	}{
		{
			RuntimeFunc: func() runtime.Runtime {
				r := runtime.Runtime{}

				r.State.Trade.Price.Last = informer.Price{Time: time.Unix(2, 0)}
				r.State.Trade.Price.Current = informer.Price{Time: time.Unix(3, 0)}
				r.Config.Trade.Pause.Min = 2 * time.Second

				return r
			},
			Expected: true,
		},
		{
			RuntimeFunc: func() runtime.Runtime {
				r := runtime.Runtime{}

				r.State.Trade.Price.Last = informer.Price{Time: time.Unix(2, 0)}
				r.State.Trade.Price.Current = informer.Price{Time: time.Unix(4, 0)}
				r.Config.Trade.Pause.Min = 2 * time.Second

				return r
			},
			Expected: false,
		},
		{
			RuntimeFunc: func() runtime.Runtime {
				r := runtime.Runtime{}

				r.State.Trade.Price.Last = informer.Price{Time: time.Unix(2, 0)}
				r.State.Trade.Price.Current = informer.Price{Time: time.Unix(5, 0)}
				r.Config.Trade.Pause.Min = 2 * time.Second

				return r
			},
			Expected: false,
		},
	}

	for i, testCase := range testCases {
		ok, err := IsInsideMinTradePause(testCase.RuntimeFunc())
		if err != nil {
			t.Fatal("case", i+1, "expected", nil, "got", err)
		}
		if ok != testCase.Expected {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", ok)
		}
	}
}
