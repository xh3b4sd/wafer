package duration

import (
	"time"
)

type Duration struct {
	// Max is the maximum time it took for the seller to emit a sell event.
	Max time.Duration
	// Min is the minimum time it took for the seller to emit a sell event.
	Min time.Duration
}
