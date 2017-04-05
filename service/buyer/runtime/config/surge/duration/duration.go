package duration

import (
	"time"
)

type Duration struct {
	// Min is the minimum time a single surge has to take place before it is
	// considered for a buy event.
	Min time.Duration `json:"min"`
}
