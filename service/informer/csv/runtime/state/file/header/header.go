package header

type Header struct {
	// Buy is the index of the row representing buy prices within the given
	// CSV file.
	Buy int
	// Ignore decides whether to ignore the first line of the given CSV.
	// This can be set to true in case the first line does not represent actual
	// data.
	Ignore bool
	// Sell is the index of the row representing sell prices within the given
	// CSV file.
	Sell int
	// Time is the index of the row representing price times within the given
	// CSV file.
	Time int
}
