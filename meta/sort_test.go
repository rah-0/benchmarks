package meta

import (
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/jfcg/sorty"
)

var sizes = []int{100, 1_000, 2_500, 5_000, 7_500, 10_000, 100_000, 1_000_000}

func BenchmarkSortInts(b *testing.B) {
	for _, size := range sizes {
		b.Run("Asc_Int_"+strconv.Itoa(size), func(b *testing.B) {
			data := genInts(size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tmp := make([]int, len(data))
				copy(tmp, data)
				sort.Ints(tmp)
			}
		})
		b.Run("Desc_Int_"+strconv.Itoa(size), func(b *testing.B) {
			data := genInts(size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tmp := make([]int, len(data))
				copy(tmp, data)
				sort.Sort(sort.Reverse(sort.IntSlice(tmp)))
			}
		})
	}
}

func BenchmarkSortyInts(b *testing.B) {
	for _, size := range sizes {
		b.Run("Asc_Int_"+strconv.Itoa(size), func(b *testing.B) {
			data := genInts(size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tmp := make([]int, len(data))
				copy(tmp, data)
				sorty.SortI(tmp)
			}
		})
	}
}

func BenchmarkSortFloats(b *testing.B) {
	for _, size := range sizes {
		b.Run("Asc_Float_"+strconv.Itoa(size), func(b *testing.B) {
			data := genFloats(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tmp := make([]float64, len(data))
				copy(tmp, data)
				sort.Float64s(tmp)
			}
		})
		b.Run("Desc_Float_"+strconv.Itoa(size), func(b *testing.B) {
			data := genFloats(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tmp := make([]float64, len(data))
				copy(tmp, data)
				sort.Sort(sort.Reverse(sort.Float64Slice(tmp)))
			}
		})
	}
}

func BenchmarkSortTimes(b *testing.B) {
	for _, size := range sizes {
		b.Run("Asc_Time_"+strconv.Itoa(size), func(b *testing.B) {
			data := genTimes(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tmp := make([]time.Time, len(data))
				copy(tmp, data)
				sort.Slice(tmp, func(i, j int) bool {
					return tmp[i].Before(tmp[j])
				})
			}
		})
		b.Run("Desc_Time_"+strconv.Itoa(size), func(b *testing.B) {
			data := genTimes(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tmp := make([]time.Time, len(data))
				copy(tmp, data)
				sort.Slice(tmp, func(i, j int) bool {
					return tmp[i].After(tmp[j])
				})
			}
		})
	}
}

func BenchmarkSortStrings(b *testing.B) {
	for _, size := range sizes {
		b.Run("Asc_String_"+strconv.Itoa(size), func(b *testing.B) {
			data := genStrings(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tmp := make([]string, len(data))
				copy(tmp, data)
				sort.Strings(tmp)
			}
		})
		b.Run("Desc_String_"+strconv.Itoa(size), func(b *testing.B) {
			data := genStrings(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tmp := make([]string, len(data))
				copy(tmp, data)
				sort.Sort(sort.Reverse(sort.StringSlice(tmp)))
			}
		})
	}
}
