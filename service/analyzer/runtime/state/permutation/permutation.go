package permutation

import (
	"time"

	"github.com/xh3b4sd/wafer/service/analyzer/runtime/state/permutation/step"
)

type Permutation struct {
	// End is the estimated time the whole permutation process is likely to come
	// to an end.
	End time.Time `json:"end"`
	// Indizes is the indizes used for the current permutational iteration of the
	// analyzer.
	Indizes  []int     `json:"indizes"`
	Max      []int     `json:"max"`
	Progress string    `json:"progress"`
	Start    time.Time `json:"start"`
	Step     step.Step `json:"step"`
}
