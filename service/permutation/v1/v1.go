package v1

import (
	"time"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	"github.com/spf13/cast"

	"github.com/xh3b4sd/wafer/service/permutation"
)

// Config is the configuration used to create a new permutation.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger
	Object permutation.Object
}

// DefaultConfig returns the default configuration used to create a new
// permutation by best effort.
func DefaultConfig() Config {
	var err error

	var newLogger micrologger.Logger
	{
		loggerConfig := micrologger.DefaultConfig()
		newLogger, err = micrologger.New(loggerConfig)
		if err != nil {
			panic(err)
		}
	}

	config := Config{
		// Dependencies.
		Logger: newLogger,
		Object: nil,
	}

	return config
}

// New creates a new configured permutation.
func New(config Config) (permutation.Permutation, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Object == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Object must not be empty")
	}

	var perms []permutation.Permutation
	{
		for _, c := range config.Object.GetPermConfigs() {
			ok := typesEqual([]interface{}{c.Min, c.Max, c.Step})
			if !ok {
				return nil, microerror.MaskAnyf(invalidConfigError, "config types must be equal")
			}

			var newPerm permutation.Permutation
			var err error

			switch t := c.Min.(type) {
			case time.Duration:
				config := DefaultDurationConfig()
				config.Min = cast.ToDuration(c.Min)
				config.Max = cast.ToDuration(c.Max)
				config.Step = cast.ToDuration(c.Step)
				newPerm, err = NewDuration(config)
				if err != nil {
					return nil, microerror.MaskAny(err)
				}
			case float64:
				config := DefaultFloat64Config()
				config.Min = cast.ToFloat64(c.Min)
				config.Max = cast.ToFloat64(c.Max)
				config.Step = cast.ToFloat64(c.Step)
				newPerm, err = NewFloat64(config)
				if err != nil {
					return nil, microerror.MaskAny(err)
				}
			default:
				return nil, microerror.MaskAnyf(invalidConfigError, "unsupported type '%T' for config value", t)
			}

			perms = append(perms, newPerm)
		}
	}

	newPermutation := &Permutation{
		// Dependencies.
		logger: config.Logger,
		object: config.Object,

		// Internals.
		max:   maxFromConfigs(config.Object.GetPermConfigs()),
		perms: perms,
	}

	return newPermutation, nil
}

// Permutation implements permutation.Permutation.
type Permutation struct {
	// Dependencies.
	logger micrologger.Logger
	object permutation.Object

	// Internals.
	max   []int
	perms []permutation.Permutation
}

func (p *Permutation) ValueFor(indizes []int) (interface{}, error) {
	if len(indizes) < len(p.max) {
		return 0, microerror.MaskAnyf(invalidExecutionError, "indizes must have length %d", len(p.max))
	}
	if indizesGreaterThan(indizes, p.max) {
		return 0, microerror.MaskAnyf(invalidExecutionError, "indizes must not be greater than max")
	}

	configs := p.object.GetPermConfigs()

	var i int
	for {
		config := configs[i]
		index := indizes[i]
		perm := p.perms[i]

		value, err := perm.ValueFor([]int{index})
		if IsInvalidExecution(err) {
			i++
			continue
		} else if err != nil {
			return nil, microerror.MaskAny(err)
		}

		err = p.object.SetPermValue(config.ID, value)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

		i++
		if i >= len(indizes) {
			break
		}
	}

	return p.object, nil
}
