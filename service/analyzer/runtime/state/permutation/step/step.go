package step

type Step struct {
	Current float64 `json:"current"`
	// Duration is the average duration a permutation step takes.
	Duration string  `json:"duration"`
	Total    float64 `json:"total"`
}
