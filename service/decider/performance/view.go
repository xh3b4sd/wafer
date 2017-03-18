package performance

import (
	"time"

	microerror "github.com/giantswarm/microkit/error"

	"github.com/xh3b4sd/wafer/service/informer"
)

type View struct {
	HasLeftNeighbour  bool
	HasRightNeighbour bool
	LeftBound         informer.Price
	RightBound        informer.Price
}

func calculateViews(histWindow []informer.Price, configView, configConvolution time.Duration) ([]View, error) {
	if len(histWindow) < 2 {
		return nil, microerror.MaskAny(notEnoughDataError)
	}

	var views []View

	views = calculateLeftBounds(histWindow, configView, configConvolution)
	views = calculateRightBounds(histWindow, views, configView)
	views = filterIncompleteViews(views, configView)
	views = setNeighbourInfo(views)

	return views, nil
}

// calculateLeftBounds sorts out left bounds of views with respect to the
// configured convolution.
func calculateLeftBounds(histWindow []informer.Price, configView, configConvolution time.Duration) []View {
	leftBound := histWindow[0].Time
	rightBound := histWindow[len(histWindow)-1].Time

	var n int
	var views []View
	for _, p := range histWindow {
		newConvolution := time.Duration(configConvolution.Seconds()*float64(n)) * time.Second
		if p.Time.Before(leftBound.Add(newConvolution)) {
			continue
		}
		view := View{
			LeftBound: p,
		}
		views = append(views, view)

		n++

		if p.Time.Add(configView).After(rightBound) {
			break
		}
	}

	return views
}

// calculateRightBounds sorts out right bounds of views according to their
// corresponding left bounds with respect to the configured view.
func calculateRightBounds(histWindow []informer.Price, views []View, configView time.Duration) []View {
	for i, v := range views {
		for _, p := range histWindow {
			if p.Time.Before(v.LeftBound.Time.Add(configView)) {
				views[i].RightBound = p
			}

			if p.Time.After(v.LeftBound.Time.Add(configView)) {
				break
			}
		}
	}

	return views
}

// filterIncompleteViews makes sure the view boundaries match with the given
// configuration.
func filterIncompleteViews(views []View, configView time.Duration) []View {
	var filteredViews []View

	tollerance := configView - time.Second

	for _, v := range views {
		if v.RightBound.Time.Before(v.LeftBound.Time.Add(tollerance)) {
			continue
		}

		filteredViews = append(filteredViews, v)
	}

	return filteredViews
}

// setNeighbourInfo makes sure each view is properly flagged with information
// about them having direct neighbour views.
func setNeighbourInfo(views []View) []View {
	for i, v := range views {
		var hasLeftNeighbour bool
		var leftNeighbour View
		if i > 0 {
			leftNeighbour = views[i-1]
			if leftNeighbour.RightBound.Time.Equal(v.LeftBound.Time) || leftNeighbour.RightBound.Time.After(v.LeftBound.Time) {
				hasLeftNeighbour = true
			}
		}

		var hasRightNeighbour bool
		var rightNeighbour View
		if i < len(views)-1 {
			rightNeighbour = views[i+1]
			if rightNeighbour.LeftBound.Time.Before(v.RightBound.Time) || rightNeighbour.LeftBound.Time.Equal(v.RightBound.Time) {
				hasRightNeighbour = true
			}
		}

		views[i].HasLeftNeighbour = hasLeftNeighbour
		views[i].HasRightNeighbour = hasRightNeighbour
	}

	return views
}
