package iteration

import (
	"io"
	"iter"
)

type WrapPullChanRowIter struct {
	next   func() (int, bool)
	stop   func()
	closed bool
}

func NewWrapPullChanIter(src IntRangeSource, gte, lt int) *WrapPullChanRowIter {
	seq := func(yield func(int) bool) {
		src.AscendRange(gte, lt, func(v int) bool { return yield(v) })
	}
	n, s := iter.Pull[int](seq)
	return &WrapPullChanRowIter{next: n, stop: s}
}

func NewWrapPullChanIterRanges(src IntRangeSource, rs []Range) *WrapPullChanRowIter {
	cp := append([]Range(nil), rs...)
	seq := func(yield func(int) bool) {
		for _, r := range cp {
			cont := true
			src.AscendRange(r.GTE, r.LT, func(v int) bool {
				if !yield(v) {
					cont = false
					return false
				}
				return true
			})
			if !cont {
				return
			}
		}
	}
	n, s := iter.Pull[int](seq)
	return &WrapPullChanRowIter{next: n, stop: s}
}

func (it *WrapPullChanRowIter) Next() (Row, error) {
	if it.closed {
		return nil, io.EOF
	}
	v, ok := it.next()
	if !ok {
		return nil, io.EOF
	}
	return Row{v}, nil
}

func (it *WrapPullChanRowIter) Close() error {
	if it.closed {
		return nil
	}
	it.closed = true

	if it.stop != nil {
		it.stop()
		for {
			_, ok := it.next()
			if !ok {
				break
			}
		}
	}

	it.next = nil
	it.stop = nil
	return nil
}
