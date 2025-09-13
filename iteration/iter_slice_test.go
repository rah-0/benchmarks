package iteration

import (
	"io"
	"reflect"
	"testing"
)

func TestSliceIter_SingleRange(t *testing.T) {
	src := &mockSource{data: []int{1, 2, 3, 4, 5}}
	it := NewSliceRowIter(src, 2, 5)

	rows, err := drainAll(t, it)
	if err != nil {
		t.Fatal(err)
	}
	got := rowsToInts(rows)
	want := []int{2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if src.calls != 1 {
		t.Fatalf("AscendRange calls=%d, want 1", src.calls)
	}
	if r, err := it.Next(); err != io.EOF || r != nil {
		t.Fatalf("Next after EOF = (%v,%v), want (nil, EOF)", r, err)
	}
}

func TestSliceIter_MultiRange(t *testing.T) {
	src := &mockSource{data: []int{0, 5, 10, 15, 20}}
	rs := []Range{{GTE: 0, LT: 1}, {GTE: 10, LT: 21}}
	it := NewSliceRowIterRanges(src, rs)

	rows, err := drainAll(t, it)
	if err != nil {
		t.Fatal(err)
	}
	got := rowsToInts(rows)
	want := []int{0, 10, 15, 20}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if src.calls != len(rs) {
		t.Fatalf("AscendRange calls=%d, want %d", src.calls, len(rs))
	}
}

func TestSliceIter_EmptyRange(t *testing.T) {
	src := &mockSource{data: []int{1, 2, 3}}
	it := NewSliceRowIter(src, 4, 5)

	if r, err := it.Next(); err != io.EOF || r != nil {
		t.Fatalf("Next = (%v,%v), want (nil, EOF)", r, err)
	}
	if src.calls != 1 {
		t.Fatalf("AscendRange calls=%d, want 1", src.calls)
	}
}

func TestSliceIter_CloseIdempotent(t *testing.T) {
	src := &mockSource{data: []int{1, 2}}
	it := NewSliceRowIter(src, 1, 3)

	if _, err := it.Next(); err != nil {
		t.Fatal(err)
	}
	if err := it.Close(); err != nil {
		t.Fatal(err)
	}
	if err := it.Close(); err != nil {
		t.Fatal(err)
	}
	if r, err := it.Next(); err != io.EOF || r != nil {
		t.Fatalf("Next after Close = (%v,%v), want (nil, EOF)", r, err)
	}
}

func TestSliceIter_ManyItems(t *testing.T) {
	N := 10_000
	data := make([]int, N)
	for i := range data {
		data[i] = i
	}
	src := &mockSource{data: data}
	it := NewSliceRowIter(src, 123, 9876)

	rows, err := drainAll(t, it)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 9876-123 {
		t.Fatalf("len(rows)=%d, want %d", len(rows), 9876-123)
	}
	if rows[0][0].(int) != 123 || rows[len(rows)-1][0].(int) != 9875 {
		t.Fatalf("bounds got [%v..%v], want [123..9875]", rows[0], rows[len(rows)-1])
	}
	if src.calls != 1 {
		t.Fatalf("AscendRange calls=%d, want 1", src.calls)
	}
}
