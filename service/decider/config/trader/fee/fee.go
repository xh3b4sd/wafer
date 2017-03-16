package fee

type Fee struct {
	// Max is the minimum amount of fees that have to be respected when
	// calculating the revenue of a single transaction.
	Max float64
	// Min is the minimum amount of fees that have to be respected when
	// calculating the revenue of a single transaction.
	Min float64
}
