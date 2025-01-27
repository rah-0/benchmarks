package testutil

func PercentDifference(big, small int) float64 {
	return (1 - float64(small)/float64(big)) * 100
}
