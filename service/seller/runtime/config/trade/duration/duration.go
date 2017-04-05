package duration

import (
	"time"
)

type Duration struct {
	// Min is the minimum time a single trade is allowed to take.
	Min time.Duration `json:"min"`
}
