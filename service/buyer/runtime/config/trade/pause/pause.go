package pause

import (
	"time"
)

type Pause struct {
	// Min is the minimum time to wait between buys.
	Min time.Duration
}
