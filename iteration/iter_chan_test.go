package iteration

import (
	"io"
	"reflect"
	"testing"
)

func TestChanIter_SingleRange(t *testing.T) {
	src := &mockSource{data: []int{1, 2, 3, 4, 5}}
	it := NewChanRowIter(src, 2, 5)

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

func TestChanIter_MultiRange(t *testing.T) {
	src := &mockSource{data: []int{0, 5, 10, 15, 20}}
	rs := []Range{{GTE: 0, LT: 1}, {GTE: 10, LT: 21}}
	it := NewChanRowIterRanges(src, rs)

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

func TestChanIter_EmptyRange(t *testing.T) {
	src := &mockSource{data: []int{1, 2, 3}}
	it := NewChanRowIter(src, 4, 5)

	if r, err := it.Next(); err != io.EOF || r != nil {
		t.Fatalf("Next = (%v,%v), want (nil, EOF)", r, err)
	}
	if src.calls != 1 {
		t.Fatalf("AscendRange calls=%d, want 1", src.calls)
	}
}

func TestChanIter_CloseIdempotent(t *testing.T) {
	src := &mockSource{data: []int{1, 2, 3}}
	it := NewChanRowIter(src, 1, 4)

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

func TestChanIter_CancelEarly(t *testing.T) {
	N := 10000
	data := make([]int, N)
	for i := 0; i < N; i++ {
		data[i] = i
	}
	src := &mockSource{data: data}
	it := NewChanRowIter(src, 0, N)

	const take = 3
	for i := 0; i < take; i++ {
		r, err := it.Next()
		if err != nil {
			t.Fatalf("Next #%d: %v", i, err)
		}
		if r[0].(int) != i {
			t.Fatalf("row=%v, want %d", r, i)
		}
	}

	if err := it.Close(); err != nil {
		t.Fatal(err)
	}

	if r, err := it.Next(); err != io.EOF || r != nil {
		t.Fatalf("Next after Close = (%v,%v), want (nil, EOF)", r, err)
	}

	if !(src.visited >= take && src.visited <= take+1) {
		t.Fatalf("visited=%d, want in [%d,%d]", src.visited, take, take+1)
	}
	if src.calls != 1 {
		t.Fatalf("AscendRange calls=%d, want 1", src.calls)
	}
}
