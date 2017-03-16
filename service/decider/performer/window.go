package performer

import (
	"time"

	"github.com/xh3b4sd/wafer/service/informer"
)

func calculateWindow(histWindow []informer.Price, configWindow time.Duration) ([]informer.Price, error) {
	// We want to ensure the chart window is always up to date. To calculate the
	// chart window we need at least 2 price events within the chart window.
	// Otherwise we cannot judge about the left and right boundaries that
	// represent the window. In case we do not have enough price events present,
	// we return the window unchanged.
	if len(histWindow) < 2 {
		return histWindow, nil
	}

	rightBound := histWindow[len(histWindow)-1].Time
	newLeftBound := rightBound.Add(-configWindow)

	var n int
	for i, p := range histWindow {
		if p.Time.Before(newLeftBound) || p.Time.Equal(newLeftBound) {
			continue
		}

		n = i
		break
	}

	newHistWindow := histWindow[n:]

	return newHistWindow, nil
}
