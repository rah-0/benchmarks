package meta

import (
	"sort"
	"strconv"
	"testing"
)

func TestParallelIntSort_SmallRandom(t *testing.T) {
	data := genInts(1000)
	expected := make([]int, len(data))
	copy(expected, data)
	sort.Ints(expected)

	ParallelIntSortAsc(data)

	if !slicesEqual(data, expected) {
		t.Errorf("sorted result incorrect for small slice")
	}
}

func TestParallelIntSort_LargeRandom(t *testing.T) {
	data := genInts(1_000_000)
	expected := make([]int, len(data))
	copy(expected, data)
	sort.Ints(expected)

	ParallelIntSortAsc(data)

	if !slicesEqual(data, expected) {
		t.Errorf("sorted result incorrect for large slice")
	}
}

// helpers
func slicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func BenchmarkParallelIntSort(b *testing.B) {
	for _, size := range sizes {
		b.Run("Asc_Int_"+strconv.Itoa(size), func(b *testing.B) {
			data := genInts(size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tmp := make([]int, len(data))
				copy(tmp, data)
				ParallelIntSortAsc(tmp)
			}
		})
		b.Run("Desc_Int_"+strconv.Itoa(size), func(b *testing.B) {
			data := genInts(size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tmp := make([]int, len(data))
				copy(tmp, data)
				ParallelIntSortDesc(tmp)
			}
		})
	}
}
