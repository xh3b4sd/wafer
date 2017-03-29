package corridor

// Corridor describes the state of a price range observed over the lifetime of
// the buyer.
type Corridor struct {
	// Max is the maximum price ever observed by the buyer.
	Max float64
}
