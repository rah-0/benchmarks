package mem

import (
	"testing"
)

func BenchmarkLargeAllocation(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slice := make([]int, 16000000)
		for j := 0; j < 16000000; j++ {
			slice[j] = j
		}
	}
}

func BenchmarkMultipleAllocations(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var slice []int
		for j := 0; j < 3000000; j++ {
			slice = append(slice, j)
		}
	}
}
