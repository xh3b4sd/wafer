package decline

type Decline struct {
	// Min is the minimum angle the observed chart is allowed to have before a
	// sell event can happen. This is used as safety net to prevent too much loss.
	Min float64
}
