package meta

import (
	"math/rand"
)

func genInts(n int) []int {
	a := make([]int, n)
	for i := range a {
		a[i] = rand.Int()
	}
	return a
}
