package duration

import (
	"time"
)

type Duration struct {
	// Max is the maximum time a single transaction is allowed to take.
	Max time.Duration
	// Min is the minimum time a single transaction is allowed to take.
	Min time.Duration
}
