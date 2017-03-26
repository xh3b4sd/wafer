package v1

import (
	"testing"
	"time"

	"github.com/spf13/cast"
)

func Test_Duration_ValueFor(t *testing.T) {
	testCases := []struct {
		Min          time.Duration
		Max          time.Duration
		Step         time.Duration
		Indizes      []int
		Expected     time.Duration
		ErrorMatcher func(err error) bool
	}{
		{
			Min:          5 * time.Second,
			Max:          30 * time.Second,
			Step:         5 * time.Second,
			Indizes:      []int{0},
			Expected:     5 * time.Second,
			ErrorMatcher: nil,
		},
		{
			Min:          5 * time.Second,
			Max:          30 * time.Second,
			Step:         5 * time.Second,
			Indizes:      []int{1},
			Expected:     10 * time.Second,
			ErrorMatcher: nil,
		},
		{
			Min:          5 * time.Second,
			Max:          30 * time.Second,
			Step:         5 * time.Second,
			Indizes:      []int{2},
			Expected:     15 * time.Second,
			ErrorMatcher: nil,
		},
		{
			Min:          5 * time.Second,
			Max:          30 * time.Second,
			Step:         5 * time.Second,
			Indizes:      []int{3},
			Expected:     20 * time.Second,
			ErrorMatcher: nil,
		},
		{
			Min:          5 * time.Second,
			Max:          30 * time.Second,
			Step:         5 * time.Second,
			Indizes:      []int{4},
			Expected:     25 * time.Second,
			ErrorMatcher: nil,
		},
		{
			Min:          5 * time.Second,
			Max:          30 * time.Second,
			Step:         5 * time.Second,
			Indizes:      []int{5},
			Expected:     30 * time.Second,
			ErrorMatcher: nil,
		},
		{
			Min:          5 * time.Second,
			Max:          30 * time.Second,
			Step:         5 * time.Second,
			Indizes:      []int{6},
			Expected:     35 * time.Second,
			ErrorMatcher: IsInvalidExecution,
		},
	}

	for i, testCase := range testCases {
		config := DefaultDurationConfig()
		config.Min = testCase.Min
		config.Max = testCase.Max
		config.Step = testCase.Step
		newDuration, err := NewDuration(config)
		if err != nil {
			t.Fatal("case", i+1, "expected", nil, "got", err)
		}

		value, err := newDuration.ValueFor(testCase.Indizes)
		if (err != nil && testCase.ErrorMatcher == nil) || (testCase.ErrorMatcher != nil && !testCase.ErrorMatcher(err)) {
			t.Fatal("case", i+1, "expected", true, "got", false)
		}

		if testCase.ErrorMatcher == nil {
			duration := cast.ToDuration(value)
			if duration != testCase.Expected {
				t.Fatal("case", i+1, "expected", testCase.Expected, "got", duration)
			}
		}
	}
}
