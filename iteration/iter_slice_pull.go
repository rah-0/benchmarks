package iteration

import (
	"io"
)

type WrapPullSliceRowIter struct {
	next   func() (int, bool)
	closed bool
}

func NewWrapPullSliceIter(src IntRangeSource, gte, lt int) *WrapPullSliceRowIter {
	buf := make([]int, 0, 64)
	src.AscendRange(gte, lt, func(v int) bool {
		buf = append(buf, v)
		return true
	})
	i := 0
	next := func() (int, bool) {
		if i >= len(buf) {
			return 0, false
		}
		v := buf[i]
		i++
		return v, true
	}
	return &WrapPullSliceRowIter{next: next}
}

func NewWrapPullSliceIterRanges(src IntRangeSource, rs []Range) *WrapPullSliceRowIter {
	cp := append([]Range(nil), rs...)
	buf := make([]int, 0, 64)
	for _, r := range cp {
		src.AscendRange(r.GTE, r.LT, func(v int) bool {
			buf = append(buf, v)
			return true
		})
	}
	i := 0
	next := func() (int, bool) {
		if i >= len(buf) {
			return 0, false
		}
		v := buf[i]
		i++
		return v, true
	}
	return &WrapPullSliceRowIter{next: next}
}

func (it *WrapPullSliceRowIter) Next() (Row, error) {
	if it.closed {
		return nil, io.EOF
	}
	v, ok := it.next()
	if !ok {
		return nil, io.EOF
	}
	return Row{v}, nil
}

func (it *WrapPullSliceRowIter) Close() error {
	if it.closed {
		return nil
	}
	it.closed = true
	it.next = nil
	return nil
}
