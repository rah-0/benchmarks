package iteration

import (
	"testing"
)

func buildSrc(n int) *mockSource {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = i
	}
	return &mockSource{data: data}
}

func burn(n int) {
	x := 0
	for i := 0; i < n; i++ {
		x += i
	}
	_ = x
}

type slowSrc struct {
	inner IntRangeSource
	work  int
}

func (s slowSrc) AscendRange(g, l int, it func(int) bool) {
	s.inner.AscendRange(g, l, func(v int) bool {
		burn(s.work)
		return it(v)
	})
}

func makeRanges(n, k int) []Range {
	rs := make([]Range, 0, k)
	step := n / k
	start := 0
	for i := 0; i < k; i++ {
		end := start + step
		if i == k-1 {
			end = n
		}
		rs = append(rs, Range{GTE: start, LT: end})
		start = end
	}
	return rs
}

func drainAllIter(it RowIter) {
	for {
		if _, err := it.Next(); err != nil {
			_ = it.Close()
			return
		}
	}
}

// ---------- FULL CONSUME

func BenchmarkPullChan_Full_100k(b *testing.B) {
	src := buildSrc(100_000)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		it := NewWrapPullChanIter(src, 0, 100_000)
		drainAllIter(it)
	}
}

func BenchmarkPullSlice_Full_100k(b *testing.B) {
	src := buildSrc(100_000)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		it := NewWrapPullSliceIter(src, 0, 100_000)
		drainAllIter(it)
	}
}

func BenchmarkBaseline_ChanDirect_Full_100k(b *testing.B) {
	src := buildSrc(100_000)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		it := NewChanRowIter(src, 0, 100_000)
		drainAllIter(it)
	}
}

func BenchmarkBaseline_SliceDirect_Full_100k(b *testing.B) {
	src := buildSrc(100_000)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		it := NewSliceRowIter(src, 0, 100_000)
		drainAllIter(it)
	}
}

// ---------- EARLY CLOSE

func benchEarlyClose(b *testing.B, ctor func(IntRangeSource, int, int) RowIter, take, n int) {
	src := buildSrc(n)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		it := ctor(src, 0, n)
		for k := 0; k < take; k++ {
			_, _ = it.Next()
		}
		_ = it.Close()
	}
}

func BenchmarkPullChan_EarlyClose10_100k(b *testing.B) {
	benchEarlyClose(b, func(s IntRangeSource, g, l int) RowIter { return NewWrapPullChanIter(s, g, l) }, 10, 100_000)
}
func BenchmarkPullSlice_EarlyClose10_100k(b *testing.B) {
	benchEarlyClose(b, func(s IntRangeSource, g, l int) RowIter { return NewWrapPullSliceIter(s, g, l) }, 10, 100_000)
}

// ---------- TIME TO FIRST ROW (TTFR)

func benchTTFR(b *testing.B, ctor func(IntRangeSource, int, int) RowIter, n int) {
	src := buildSrc(n)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		it := ctor(src, 0, n)
		_, _ = it.Next()
		_ = it.Close()
	}
}

func BenchmarkPullChan_TTFR_1M(b *testing.B) {
	benchTTFR(b, func(s IntRangeSource, g, l int) RowIter { return NewWrapPullChanIter(s, g, l) }, 1_000_000)
}
func BenchmarkPullSlice_TTFR_1M(b *testing.B) {
	benchTTFR(b, func(s IntRangeSource, g, l int) RowIter { return NewWrapPullSliceIter(s, g, l) }, 1_000_000)
}

// ---------- SLOW CONSUMER

func benchSlowConsumer(b *testing.B, ctor func(IntRangeSource, int, int) RowIter, n, work int) {
	src := buildSrc(n)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		it := ctor(src, 0, n)
		for {
			_, err := it.Next()
			if err != nil {
				break
			}
			burn(work)
		}
		_ = it.Close()
	}
}

func BenchmarkPullChan_SlowConsumerWork50_100k(b *testing.B) {
	benchSlowConsumer(b, func(s IntRangeSource, g, l int) RowIter { return NewWrapPullChanIter(s, g, l) }, 100_000, 50)
}
func BenchmarkPullSlice_SlowConsumerWork50_100k(b *testing.B) {
	benchSlowConsumer(b, func(s IntRangeSource, g, l int) RowIter { return NewWrapPullSliceIter(s, g, l) }, 100_000, 50)
}

// ---------- SLOW PRODUCER

func benchSlowProducer(b *testing.B, ctor func(IntRangeSource, int, int) RowIter, n, work int) {
	base := buildSrc(n)
	src := slowSrc{inner: base, work: work}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		it := ctor(src, 0, n)
		drainAllIter(it)
	}
}

func BenchmarkPullChan_SlowProducerWork50_100k(b *testing.B) {
	benchSlowProducer(b, func(s IntRangeSource, g, l int) RowIter { return NewWrapPullChanIter(s, g, l) }, 100_000, 50)
}
func BenchmarkPullSlice_SlowProducerWork50_100k(b *testing.B) {
	benchSlowProducer(b, func(s IntRangeSource, g, l int) RowIter { return NewWrapPullSliceIter(s, g, l) }, 100_000, 50)
}

// ---------- MULTI-RANGE VS SINGLE

func benchMultiRange(b *testing.B, ctorR func(IntRangeSource, []Range) RowIter, n, parts int) {
	src := buildSrc(n)
	rs := makeRanges(n, parts)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		it := ctorR(src, rs)
		drainAllIter(it)
	}
}

func BenchmarkPullChan_MultiRange_100k_10parts(b *testing.B) {
	benchMultiRange(b, func(s IntRangeSource, rs []Range) RowIter { return NewWrapPullChanIterRanges(s, rs) }, 100_000, 10)
}
func BenchmarkPullSlice_MultiRange_100k_10parts(b *testing.B) {
	benchMultiRange(b, func(s IntRangeSource, rs []Range) RowIter { return NewWrapPullSliceIterRanges(s, rs) }, 100_000, 10)
}

// ---------- EMPTY RANGES

func BenchmarkPullChan_EmptyRange(b *testing.B) {
	src := buildSrc(0)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		it := NewWrapPullChanIter(src, 0, 0)
		drainAllIter(it)
	}
}

func BenchmarkPullSlice_EmptyRange(b *testing.B) {
	src := buildSrc(0)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		it := NewWrapPullSliceIter(src, 0, 0)
		drainAllIter(it)
	}
}

// ---------- PARALLEL (many iterators at once)

func BenchmarkPullChan_Parallel_10k(b *testing.B) {
	const N = 10_000
	src := buildSrc(N)
	b.ReportAllocs()
	b.SetParallelism(4)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			it := NewWrapPullChanIter(src, 0, N)
			drainAllIter(it)
		}
	})
}

func BenchmarkPullSlice_Parallel_10k(b *testing.B) {
	const N = 10_000
	src := buildSrc(N)
	b.ReportAllocs()
	b.SetParallelism(4)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			it := NewWrapPullSliceIter(src, 0, N)
			drainAllIter(it)
		}
	})
}
