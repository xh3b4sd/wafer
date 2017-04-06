package price

import (
	"time"
)

type Price struct {
	End    time.Time `json:"end"`
	Events int       `json:"events"`
	Start  time.Time `json:"start"`
}
