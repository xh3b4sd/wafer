package v1

import (
	"time"
)

const (
	DefaultDurationCap = 10
)

type Duration struct {
	timeDurations []time.Duration
}

func (d *Duration) Add(td time.Duration) {
	d.timeDurations = append(d.timeDurations, td)

	if len(d.timeDurations) > DefaultDurationCap {
		d.timeDurations = d.timeDurations[1:]
	}
}

func (d *Duration) Average() time.Duration {
	amount := int64(len(d.timeDurations))
	var total int64

	for _, td := range d.timeDurations {
		total += td.Nanoseconds()
	}

	average := time.Duration(total/amount) * time.Nanosecond

	return average
}
