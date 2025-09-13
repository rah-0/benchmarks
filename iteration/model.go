package iteration

import (
	"io"
	"testing"
)

type Row = []any

type RowIter interface {
	Next() (Row, error)
	Close() error
}

type Range struct{ GTE, LT int }

type IntRangeSource interface {
	AscendRange(gte, lt int, it func(int) bool)
}

type mockSource struct {
	data    []int
	calls   int
	visited int
}

func (s *mockSource) AscendRange(gte, lt int, it func(int) bool) {
	s.calls++
	for _, v := range s.data {
		if v < gte || v >= lt {
			continue
		}
		s.visited++
		if !it(v) {
			return
		}
	}
}

func rowsToInts(rows []Row) []int {
	out := make([]int, len(rows))
	for i, r := range rows {
		out[i] = r[0].(int)
	}
	return out
}

func drainAll(t *testing.T, it RowIter) ([]Row, error) {
	t.Helper()
	var out []Row
	for {
		r, err := it.Next()
		if err == io.EOF {
			return out, nil
		}
		if err != nil {
			return out, err
		}
		out = append(out, r)
	}
}
