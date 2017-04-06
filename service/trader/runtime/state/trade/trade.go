package trade

type Trade struct {
	// Cycles is the number of buy and sell iterations the trader processed so
	// far. After one buy must come one sell. Each item of the list represents the
	// result of the informer's price events lists.
	Cycles []int64
	// Revenue is the total amount of revenue the seller made so far. Each item of
	// the list represents the result of the informer's price events lists.
	Revenues []float64
}
