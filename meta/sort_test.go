package meta

import (
	"sort"
	"strconv"
	"testing"

	pargosort "github.com/exascience/pargo/sort"
	"github.com/jfcg/sorty"
	"github.com/rah-0/parsort"
)

var testSizes = []int{10_000, 100_000, 1_000_000, 10_000_000}

func BenchmarkSortInts(b *testing.B) {
	for _, size := range testSizes {
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
	}
}

func BenchmarkSortyInts(b *testing.B) {
	for _, size := range testSizes {
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

type pargoInts []int

func (p pargoInts) Len() int           { return len(p) }
func (p pargoInts) Less(i, j int) bool { return p[i] < p[j] }
func (p pargoInts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pargoInts) SequentialSort(i, j int) {
	sort.Slice(p[i:j], func(x, y int) bool {
		return p[i+x] < p[i+y]
	})
}

func BenchmarkPargoInts(b *testing.B) {
	for _, size := range testSizes {
		b.Run("Asc_Int_"+strconv.Itoa(size), func(b *testing.B) {
			data := genInts(size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tmp := make([]int, len(data))
				copy(tmp, data)
				pargosort.Sort(pargoInts(tmp))
			}
		})
	}
}

func BenchmarkParsortInts(b *testing.B) {
	for _, size := range testSizes {
		b.Run("Asc_Int_"+strconv.Itoa(size), func(b *testing.B) {
			data := genInts(size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tmp := make([]int, len(data))
				copy(tmp, data)
				parsort.IntAsc(tmp)
			}
		})
	}
}
