package v1

import (
	"reflect"
	"testing"
	"time"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/spf13/cast"

	"github.com/xh3b4sd/wafer/service/permutation"
	permutationconfig "github.com/xh3b4sd/wafer/service/permutation/runtime/config"
)

func Test_Permutation_ValueFor(t *testing.T) {
	testCases := []struct {
		Indizes      []int
		Expected     *testConfig
		ErrorMatcher func(err error) bool
	}{
		{
			Indizes: []int{0, 0, 0},
			Expected: &testConfig{
				Foo: 0 * time.Second,
				Bar: float64(0.05),
				Baz: 0 * time.Second,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{1, 0, 0},
			Expected: &testConfig{
				Foo: 2 * time.Second,
				Bar: float64(0.05),
				Baz: 0 * time.Second,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{2, 0, 0},
			Expected: &testConfig{
				Foo: 4 * time.Second,
				Bar: float64(0.05),
				Baz: 0 * time.Second,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{3, 0, 0},
			Expected: &testConfig{
				Foo: 6 * time.Second,
				Bar: float64(0.05),
				Baz: 0 * time.Second,
			},
			ErrorMatcher: nil,
		},
		// Test case 5.
		{
			Indizes: []int{0, 1, 0},
			Expected: &testConfig{
				Foo: 0 * time.Second,
				Bar: float64(0.10),
				Baz: 0 * time.Second,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{1, 1, 0},
			Expected: &testConfig{
				Foo: 2 * time.Second,
				Bar: float64(0.10),
				Baz: 0 * time.Second,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{2, 1, 0},
			Expected: &testConfig{
				Foo: 4 * time.Second,
				Bar: float64(0.10),
				Baz: 0 * time.Second,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{3, 1, 0},
			Expected: &testConfig{
				Foo: 6 * time.Second,
				Bar: float64(0.10),
				Baz: 0 * time.Second,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{0, 0, 1},
			Expected: &testConfig{
				Foo: 0 * time.Second,
				Bar: float64(0.05),
				Baz: 10 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		// Test case 10.
		{
			Indizes: []int{1, 0, 1},
			Expected: &testConfig{
				Foo: 2 * time.Second,
				Bar: float64(0.05),
				Baz: 10 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{2, 0, 1},
			Expected: &testConfig{
				Foo: 4 * time.Second,
				Bar: float64(0.05),
				Baz: 10 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{3, 0, 1},
			Expected: &testConfig{
				Foo: 6 * time.Second,
				Bar: float64(0.05),
				Baz: 10 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{0, 1, 1},
			Expected: &testConfig{
				Foo: 0 * time.Second,
				Bar: float64(0.10),
				Baz: 10 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{1, 1, 1},
			Expected: &testConfig{
				Foo: 2 * time.Second,
				Bar: float64(0.10),
				Baz: 10 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		// Test case 15.
		{
			Indizes: []int{2, 1, 1},
			Expected: &testConfig{
				Foo: 4 * time.Second,
				Bar: float64(0.10),
				Baz: 10 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{3, 1, 1},
			Expected: &testConfig{
				Foo: 6 * time.Second,
				Bar: float64(0.10),
				Baz: 10 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{0, 0, 2},
			Expected: &testConfig{
				Foo: 0 * time.Second,
				Bar: float64(0.05),
				Baz: 20 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{1, 0, 2},
			Expected: &testConfig{
				Foo: 2 * time.Second,
				Bar: float64(0.05),
				Baz: 20 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{2, 0, 2},
			Expected: &testConfig{
				Foo: 4 * time.Second,
				Bar: float64(0.05),
				Baz: 20 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		// Test case 20.
		{
			Indizes: []int{3, 0, 2},
			Expected: &testConfig{
				Foo: 6 * time.Second,
				Bar: float64(0.05),
				Baz: 20 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{0, 1, 2},
			Expected: &testConfig{
				Foo: 0 * time.Second,
				Bar: float64(0.10),
				Baz: 20 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{1, 1, 2},
			Expected: &testConfig{
				Foo: 2 * time.Second,
				Bar: float64(0.10),
				Baz: 20 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{2, 1, 2},
			Expected: &testConfig{
				Foo: 4 * time.Second,
				Bar: float64(0.10),
				Baz: 20 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		{
			Indizes: []int{3, 1, 2},
			Expected: &testConfig{
				Foo: 6 * time.Second,
				Bar: float64(0.10),
				Baz: 20 * time.Hour,
			},
			ErrorMatcher: nil,
		},
		// Test case 25.
		{
			Indizes:      []int{4, 1, 2},
			Expected:     nil,
			ErrorMatcher: IsInvalidExecution,
		},
		{
			Indizes:      []int{3, 2, 2},
			Expected:     nil,
			ErrorMatcher: IsInvalidExecution,
		},
		{
			Indizes:      []int{3, 1, 3},
			Expected:     nil,
			ErrorMatcher: IsInvalidExecution,
		},
		{
			Indizes:      []int{3, 1, 2, 5},
			Expected:     nil,
			ErrorMatcher: IsInvalidExecution,
		},
		{
			Indizes:      []int{4, 2, 3},
			Expected:     nil,
			ErrorMatcher: IsInvalidExecution,
		},
		// Test case 30.
		{
			Indizes:      []int{4, 2, 3, 5},
			Expected:     nil,
			ErrorMatcher: IsInvalidExecution,
		},
		{
			Indizes:      []int{},
			Expected:     nil,
			ErrorMatcher: IsInvalidExecution,
		},
		{
			Indizes:      []int{3},
			Expected:     nil,
			ErrorMatcher: IsInvalidExecution,
		},
		{
			Indizes:      []int{3, 1},
			Expected:     nil,
			ErrorMatcher: IsInvalidExecution,
		},
	}

	var err error

	var newPermutation permutation.Permutation
	{
		config := DefaultConfig()
		config.Object = &testConfig{}
		newPermutation, err = New(config)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	for i, testCase := range testCases {
		if i != 24 {
			continue
		}
		value, err := newPermutation.ValueFor(testCase.Indizes)
		if (err != nil && testCase.ErrorMatcher == nil) || (testCase.ErrorMatcher != nil && !testCase.ErrorMatcher(err)) {
			t.Fatal("case", i+1, "expected", true, "got", false)
		}

		if testCase.ErrorMatcher == nil {
			object, ok := value.(*testConfig)
			if !ok {
				t.Fatal("case", i+1, "expected", true, "got", false)
			}

			if !reflect.DeepEqual(object, testCase.Expected) {
				t.Fatal("case", i+1, "expected", testCase.Expected, "got", object)
			}
		}
	}
}

const (
	testPermIDFoo = "Foo"
	testPermIDBar = "Bar"
	testPermIDBaz = "Baz"
)

type testConfig struct {
	Foo time.Duration
	Bar float64
	Baz time.Duration
}

func (c *testConfig) GetPermConfigs() []permutationconfig.Config {
	var config permutationconfig.Config
	var configs []permutationconfig.Config

	config = permutationconfig.Config{}
	config.ID = testPermIDFoo
	config.Min = 0 * time.Second
	config.Max = 6 * time.Second
	config.Step = 2 * time.Second
	configs = append(configs, config)

	config = permutationconfig.Config{}
	config.ID = testPermIDBar
	config.Min = 0.05
	config.Max = 0.10
	config.Step = 0.05
	configs = append(configs, config)

	config = permutationconfig.Config{}
	config.ID = testPermIDBaz
	config.Min = 0 * time.Second
	config.Max = 20 * time.Hour
	config.Step = 10 * time.Hour
	configs = append(configs, config)

	return configs
}

func (c *testConfig) SetPermValue(permID string, permValue interface{}) error {
	switch permID {
	case testPermIDFoo:
		d, err := cast.ToDurationE(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Foo = d
	case testPermIDBar:
		f, err := cast.ToFloat64E(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Bar = f
	case testPermIDBaz:
		d, err := cast.ToDurationE(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Baz = d
	default:
		return microerror.MaskAnyf(invalidExecutionError, "unknown permID '%s'", permID)
	}

	return nil
}
