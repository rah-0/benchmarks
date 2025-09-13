package iteration

import (
	"io"
	"reflect"
	"testing"
	"time"
)

func TestPullIter_SingleRange(t *testing.T) {
	src := &mockSource{data: []int{1, 2, 3, 4, 5}}
	it := NewWrapPullChanIter(src, 2, 5)

	rows, err := drainAll(t, it)
	if err != nil {
		t.Fatal(err)
	}
	got := rowsToInts(rows)
	want := []int{2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if r, err := it.Next(); err != io.EOF || r != nil {
		t.Fatalf("Next after EOF = (%v,%v), want (nil, EOF)", r, err)
	}
}

func TestPullIter_MultiRange(t *testing.T) {
	src := &mockSource{data: []int{0, 5, 10, 15, 20}}
	rs := []Range{{GTE: 0, LT: 1}, {GTE: 10, LT: 21}}
	it := NewWrapPullChanIterRanges(src, rs)

	rows, err := drainAll(t, it)
	if err != nil {
		t.Fatal(err)
	}
	got := rowsToInts(rows)
	want := []int{0, 10, 15, 20}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestPullIter_EmptyRange(t *testing.T) {
	src := &mockSource{data: []int{1, 2, 3}}
	it := NewWrapPullChanIter(src, 4, 5)

	if r, err := it.Next(); err != io.EOF || r != nil {
		t.Fatalf("Next = (%v,%v), want (nil, EOF)", r, err)
	}
}

func TestPullIter_CloseIdempotent(t *testing.T) {
	src := &mockSource{data: []int{1, 2, 3, 4}}
	it := NewWrapPullChanIter(src, 1, 4)

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

func TestPullIter_CancelEarly(t *testing.T) {
	N := 10000
	data := make([]int, N)
	for i := 0; i < N; i++ {
		data[i] = i
	}
	src := &mockSource{data: data}
	it := NewWrapPullChanIter(src, 0, N)

	for i := 0; i < 3; i++ {
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

	done := make(chan struct{})
	go func() {
		_, _ = it.Next()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("Next after Close did not return promptly")
	}
}
