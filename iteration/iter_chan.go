package iteration

import (
	"io"
	"sync"
)

type ChanRowIter struct {
	src      IntRangeSource
	ranges   []Range
	rangeIdx int

	ch     chan *Row
	done   chan struct{}
	wg     sync.WaitGroup
	closed bool
}

func NewChanRowIter(src IntRangeSource, gte, lt int) *ChanRowIter {
	return &ChanRowIter{
		src:      src,
		ranges:   []Range{{GTE: gte, LT: lt}},
		rangeIdx: -1,
		done:     make(chan struct{}),
	}
}

func NewChanRowIterRanges(src IntRangeSource, rs []Range) *ChanRowIter {
	cp := append([]Range(nil), rs...)
	return &ChanRowIter{
		src:      src,
		ranges:   cp,
		rangeIdx: -1,
		done:     make(chan struct{}),
	}
}

func (it *ChanRowIter) startNextRange() bool {
	it.rangeIdx++
	if it.rangeIdx >= len(it.ranges) {
		return false
	}
	r := it.ranges[it.rangeIdx]
	it.ch = make(chan *Row)

	it.wg.Add(1)
	go func(gte, lt int, out chan<- *Row, stop <-chan struct{}) {
		defer it.wg.Done()
		defer close(out)
		it.src.AscendRange(gte, lt, func(v int) bool {
			row := Row{v}
			select {
			case out <- &row:
				return true
			case <-stop:
				return false
			}
		})
	}(r.GTE, r.LT, it.ch, it.done)

	return true
}

func (it *ChanRowIter) Next() (Row, error) {
	if it.closed {
		return nil, io.EOF
	}
	for it.ch == nil {
		if !it.startNextRange() {
			return nil, io.EOF
		}
	}
	ptr, ok := <-it.ch
	if ok && ptr != nil {
		return *ptr, nil
	}
	it.ch = nil
	return it.Next()
}

func (it *ChanRowIter) Close() error {
	if it.closed {
		return nil
	}
	it.closed = true
	close(it.done)
	it.wg.Wait()
	it.ch = nil
	it.ranges = nil
	it.src = nil
	return nil
}
