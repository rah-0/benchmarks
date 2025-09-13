package iteration

import (
	"io"
)

type SliceRowIter struct {
	src    IntRangeSource
	ranges []Range
	buf    []int
	index  int
	filled bool
	closed bool
}

func NewSliceRowIter(src IntRangeSource, gte, lt int) *SliceRowIter {
	return &SliceRowIter{
		src:    src,
		ranges: []Range{{GTE: gte, LT: lt}},
	}
}

func NewSliceRowIterRanges(src IntRangeSource, rs []Range) *SliceRowIter {
	cp := append([]Range(nil), rs...)
	return &SliceRowIter{
		src:    src,
		ranges: cp,
	}
}

func (it *SliceRowIter) fillOnce() {
	if it.filled {
		return
	}
	it.buf = make([]int, 0, 64)
	for _, r := range it.ranges {
		it.src.AscendRange(r.GTE, r.LT, func(v int) bool {
			it.buf = append(it.buf, v)
			return true
		})
	}
	it.filled = true
	it.index = 0
}

func (it *SliceRowIter) Next() (Row, error) {
	if it.closed {
		return nil, io.EOF
	}
	it.fillOnce()
	if it.index >= len(it.buf) {
		return nil, io.EOF
	}
	v := it.buf[it.index]
	it.index++
	return Row{v}, nil
}

func (it *SliceRowIter) Close() error {
	if it.closed {
		return nil
	}
	it.closed = true
	it.buf = nil
	it.ranges = nil
	it.src = nil
	return nil
}
