package meta

import (
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// sinkInt64 prevents the compiler from eliminating benchmark calls whose
// return values would otherwise be unused (dead-code elimination).
var sinkInt64 int64

func TestAtomicCounter(t *testing.T) {
	var c AtomicCounter
	for i := 0; i < 1000; i++ {
		c.Increment()
	}
	assert.Equal(t, int64(1000), c.Get())
}

func TestMutexCounter(t *testing.T) {
	var c MutexCounter
	for i := 0; i < 1000; i++ {
		c.Increment()
	}
	assert.Equal(t, int64(1000), c.Get())
}

func TestRWMutexCounter(t *testing.T) {
	var c RWMutexCounter
	for i := 0; i < 1000; i++ {
		c.Increment()
	}
	assert.Equal(t, int64(1000), c.Get())
}

func TestChannelCounter(t *testing.T) {
	c := NewChannelCounter(256)
	defer c.Close()
	for i := 0; i < 1000; i++ {
		c.Increment()
	}
	reply := make(chan int64, 1)
	assert.Equal(t, int64(1000), c.Get(reply))
}

func TestCASSpinCounter(t *testing.T) {
	var c CASSpinCounter
	for i := 0; i < 1000; i++ {
		c.Increment()
	}
	assert.Equal(t, int64(1000), c.Get())
}

func TestCASBackoffCounter(t *testing.T) {
	var c CASBackoffCounter
	for i := 0; i < 1000; i++ {
		c.Increment()
	}
	assert.Equal(t, int64(1000), c.Get())
}

func TestAtomicCounter_Concurrent(t *testing.T) {
	var c AtomicCounter
	runConcurrentIncrements(t, 100, 1000, c.Increment)
	assert.Equal(t, int64(100*1000), c.Get())
}

func TestMutexCounter_Concurrent(t *testing.T) {
	var c MutexCounter
	runConcurrentIncrements(t, 100, 1000, c.Increment)
	assert.Equal(t, int64(100*1000), c.Get())
}

func TestRWMutexCounter_Concurrent(t *testing.T) {
	var c RWMutexCounter
	runConcurrentIncrements(t, 100, 1000, c.Increment)
	assert.Equal(t, int64(100*1000), c.Get())
}

func TestChannelCounter_Concurrent(t *testing.T) {
	c := NewChannelCounter(256)
	defer c.Close()
	runConcurrentIncrements(t, 100, 1000, c.Increment)
	reply := make(chan int64, 1)
	assert.Equal(t, int64(100*1000), c.Get(reply))
}

func TestCASSpinCounter_Concurrent(t *testing.T) {
	var c CASSpinCounter
	runConcurrentIncrements(t, 100, 1000, c.Increment)
	assert.Equal(t, int64(100*1000), c.Get())
}

func TestCASBackoffCounter_Concurrent(t *testing.T) {
	var c CASBackoffCounter
	runConcurrentIncrements(t, 100, 1000, c.Increment)
	assert.Equal(t, int64(100*1000), c.Get())
}

func runConcurrentIncrements(t *testing.T, goroutines, incsPerGoroutine int, inc func()) {
	t.Helper()
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < incsPerGoroutine; i++ {
				inc()
			}
		}()
	}
	wg.Wait()
}

func TestSPSCRingBuffer_Serial(t *testing.T) {
	rb := NewSPSCRingBuffer(8)

	// Fill it up (capacity is 7 usable slots for power-of-two ring)
	for i := int64(0); i < 7; i++ {
		assert.True(t, rb.Push(i), "Push %d should succeed", i)
	}
	assert.False(t, rb.Push(99), "Push to full buffer should fail")

	// Drain and verify order
	for i := int64(0); i < 7; i++ {
		val, ok := rb.Pop()
		assert.True(t, ok)
		assert.Equal(t, i, val)
	}
	_, ok := rb.Pop()
	assert.False(t, ok, "Pop from empty buffer should fail")
}

func TestSPSCRingBuffer_ProducerConsumer(t *testing.T) {
	const count = 100_000
	rb := NewSPSCRingBuffer(1024)

	done := make(chan struct{})
	go func() {
		for i := int64(0); i < count; i++ {
			for !rb.Push(i) {
				runtime.Gosched()
			}
		}
		close(done)
	}()

	for i := int64(0); i < count; i++ {
		var val int64
		var ok bool
		for {
			val, ok = rb.Pop()
			if ok {
				break
			}
			runtime.Gosched()
		}
		assert.Equal(t, i, val)
	}
	<-done
}

func TestSyncMapStore(t *testing.T) {
	var s SyncMapStore
	s.Store("a", 1)
	s.Store("b", 2)

	v, ok := s.Load("a")
	assert.True(t, ok)
	assert.Equal(t, int64(1), v)

	v, ok = s.Load("b")
	assert.True(t, ok)
	assert.Equal(t, int64(2), v)

	_, ok = s.Load("missing")
	assert.False(t, ok)
}

func TestMutexMap(t *testing.T) {
	m := NewMutexMap()
	m.Store("a", 1)
	m.Store("b", 2)

	v, ok := m.Load("a")
	assert.True(t, ok)
	assert.Equal(t, int64(1), v)

	v, ok = m.Load("b")
	assert.True(t, ok)
	assert.Equal(t, int64(2), v)

	_, ok = m.Load("missing")
	assert.False(t, ok)
}

func TestSemaphore(t *testing.T) {
	sem := NewSemaphore(2)

	// Should be able to acquire twice without blocking
	sem.Acquire()
	sem.Acquire()

	// Release one, acquire again
	sem.Release()
	sem.Acquire()

	sem.Release()
	sem.Release()
}

func TestWaitGroupWork(t *testing.T) {
	var counter AtomicCounter
	WaitGroupWork(10, func() {
		counter.Increment()
	})
	assert.Equal(t, int64(10), counter.Get())
}

func TestPoolGetPut(t *testing.T) {
	// Just verify it doesn't panic
	for i := 0; i < 100; i++ {
		PoolGetPut()
	}
}

func TestOnceInit(t *testing.T) {
	var o OnceInit
	assert.Equal(t, int64(0), o.Get())
	o.Init()
	assert.Equal(t, int64(42), o.Get())
	// Calling Init again should be a no-op
	o.Init()
	assert.Equal(t, int64(42), o.Get())
}

func TestOnceInit_Concurrent(t *testing.T) {
	var o OnceInit
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			o.Init()
		}()
	}
	wg.Wait()
	assert.Equal(t, int64(42), o.Get())
}

// ---------------------------------------------------------------------------
// Sharded Counter correctness
// ---------------------------------------------------------------------------

func TestShardedCounter(t *testing.T) {
	c := NewShardedCounter(8)
	for i := 0; i < 1000; i++ {
		c.Increment(i)
	}
	assert.Equal(t, int64(1000), c.Get())
}

func TestShardedCounter_Concurrent(t *testing.T) {
	c := NewShardedCounter(8)
	var wg sync.WaitGroup
	wg.Add(100)
	for g := 0; g < 100; g++ {
		g := g
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				c.Increment(g)
			}
		}()
	}
	wg.Wait()
	assert.Equal(t, int64(100*1000), c.Get())
}

// ---------------------------------------------------------------------------
// Per-CPU Counter correctness
// ---------------------------------------------------------------------------

func TestPerCPUCounter(t *testing.T) {
	c := NewPerCPUCounter()
	for i := 0; i < 1000; i++ {
		c.Increment(i)
	}
	assert.Equal(t, int64(1000), c.Get())
}

func TestPerCPUCounter_Concurrent(t *testing.T) {
	c := NewPerCPUCounter()
	var wg sync.WaitGroup
	wg.Add(100)
	for g := 0; g < 100; g++ {
		g := g
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				c.Increment(g)
			}
		}()
	}
	wg.Wait()
	assert.Equal(t, int64(100*1000), c.Get())
}

// ---------------------------------------------------------------------------
// Flat Combining Counter correctness
// ---------------------------------------------------------------------------

func TestFlatCombiningCounter(t *testing.T) {
	var c FlatCombiningCounter
	for i := 0; i < 1000; i++ {
		c.Increment()
	}
	assert.Equal(t, int64(1000), c.Get())
}

func TestFlatCombiningCounter_Concurrent(t *testing.T) {
	var c FlatCombiningCounter
	runConcurrentIncrements(t, 100, 1000, c.Increment)
	assert.Equal(t, int64(100*1000), c.Get())
}

// ---------------------------------------------------------------------------
// Local Buffered Counter correctness
// ---------------------------------------------------------------------------

func TestLocalBufferedCounter(t *testing.T) {
	c := NewLocalBufferedCounter(64)
	var local int64
	for i := 0; i < 1000; i++ {
		c.IncrementBuffered(&local)
	}
	c.Flush(&local)
	assert.Equal(t, int64(1000), c.Get())
}

func TestLocalBufferedCounter_Concurrent(t *testing.T) {
	c := NewLocalBufferedCounter(64)
	var wg sync.WaitGroup
	wg.Add(100)
	for g := 0; g < 100; g++ {
		go func() {
			defer wg.Done()
			var local int64
			for i := 0; i < 1000; i++ {
				c.IncrementBuffered(&local)
			}
			c.Flush(&local)
		}()
	}
	wg.Wait()
	assert.Equal(t, int64(100*1000), c.Get())
}

// ---------------------------------------------------------------------------
// MPSC Aggregator correctness
// ---------------------------------------------------------------------------

func TestMPSCAggregator(t *testing.T) {
	a := NewMPSCAggregator(256)
	defer a.Close()
	for i := 0; i < 1000; i++ {
		a.Increment()
	}
	// Give aggregator time to drain
	var val int64
	for attempts := 0; attempts < 200; attempts++ {
		val = a.Get()
		if val == 1000 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	assert.Equal(t, int64(1000), val)
}

func TestMPSCAggregator_Concurrent(t *testing.T) {
	a := NewMPSCAggregator(4096)
	defer a.Close()
	var wg sync.WaitGroup
	wg.Add(100)
	for g := 0; g < 100; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				a.Increment()
			}
		}()
	}
	wg.Wait()
	// Give aggregator goroutine time to drain the channel
	var val int64
	for attempts := 0; attempts < 200; attempts++ {
		val = a.Get()
		if val == 100*1000 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	assert.Equal(t, int64(100*1000), val)
}

// ---------------------------------------------------------------------------
// RCU Counter correctness
// ---------------------------------------------------------------------------

func TestRCUCounter(t *testing.T) {
	c := NewRCUCounter()
	for i := 0; i < 1000; i++ {
		c.Increment()
	}
	assert.Equal(t, int64(1000), c.Get())
}

func TestRCUCounter_Concurrent(t *testing.T) {
	c := NewRCUCounter()
	runConcurrentIncrements(t, 100, 1000, c.Increment)
	assert.Equal(t, int64(100*1000), c.Get())
}

// ---------------------------------------------------------------------------
// Disruptor MPSC correctness
// ---------------------------------------------------------------------------

func TestDisruptorMPSC_Serial(t *testing.T) {
	d := NewDisruptorMPSC(8)
	for i := int64(0); i < 7; i++ {
		assert.True(t, d.Publish(i+1), "Publish %d should succeed", i+1)
	}
	// Buffer should be full (8 slots, but sequence-based so 8 usable)
	// Actually with sequence math wc-rc > mask means 8 items fit
	for i := int64(0); i < 7; i++ {
		val, ok := d.Consume()
		assert.True(t, ok)
		assert.Equal(t, i+1, val)
	}
	_, ok := d.Consume()
	assert.False(t, ok, "Consume from empty should fail")
}

func TestDisruptorMPSC_ProducerConsumer(t *testing.T) {
	const count = 10_000
	d := NewDisruptorMPSC(1024)

	var wg sync.WaitGroup
	wg.Add(4) // 4 producers
	for p := 0; p < 4; p++ {
		go func() {
			defer wg.Done()
			for i := 0; i < count/4; i++ {
				for !d.Publish(1) {
					runtime.Gosched()
				}
			}
		}()
	}

	var total int64
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	finished := false
	for !finished {
		select {
		case <-done:
			finished = true
		default:
		}
		if val, ok := d.Consume(); ok {
			total += val
		} else {
			runtime.Gosched()
		}
	}
	// Drain remaining
	for {
		val, ok := d.Consume()
		if !ok {
			break
		}
		total += val
	}
	assert.Equal(t, int64(count), total)
}

// ---------------------------------------------------------------------------
// Approximate Counter correctness
// ---------------------------------------------------------------------------

func TestApproxCounter(t *testing.T) {
	c := NewApproxCounter(1) // sampleRate=1 means every increment counts
	seed := uint64(12345)
	for i := 0; i < 10000; i++ {
		c.Increment(&seed)
	}
	assert.Equal(t, int64(10000), c.Get())
}

func TestApproxCounter_Approximate(t *testing.T) {
	c := NewApproxCounter(100)
	seed := uint64(67890)
	n := 1_000_000
	for i := 0; i < n; i++ {
		c.Increment(&seed)
	}
	got := c.Get()
	// With sampleRate=100, expect ~1M but with statistical variance
	// Allow 20% tolerance
	assert.InDelta(t, float64(n), float64(got), float64(n)*0.2,
		"expected ~%d, got %d", n, got)
}


func BenchmarkCounters_IncOnly(b *testing.B) {
	b.Run("Atomic", func(b *testing.B) {
		var c AtomicCounter
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				c.Increment()
			}
		})
	})

	b.Run("Mutex", func(b *testing.B) {
		var c MutexCounter
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				c.Increment()
			}
		})
	})

	b.Run("RWMutex", func(b *testing.B) {
		var c RWMutexCounter
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				c.Increment()
			}
		})
	})

	b.Run("CAS_Spin", func(b *testing.B) {
		var c CASSpinCounter
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				c.Increment()
			}
		})
	})

	b.Run("CAS_Backoff", func(b *testing.B) {
		var c CASBackoffCounter
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				c.Increment()
			}
		})
	})

	for _, buf := range []int{0, 1, 256, 4096} {
		b.Run("Channel_buf"+strconv.Itoa(buf), func(b *testing.B) {
			c := NewChannelCounter(buf)
			b.Cleanup(c.Close)
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					c.Increment()
				}
			})
		})
	}
}

func BenchmarkCounters_GetOnly(b *testing.B) {
	b.Run("Atomic", func(b *testing.B) {
		var c AtomicCounter
		c.Increment() // seed
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				sinkInt64 = c.Get()
			}
		})
	})

	b.Run("Mutex", func(b *testing.B) {
		var c MutexCounter
		c.Increment()
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				sinkInt64 = c.Get()
			}
		})
	})

	b.Run("RWMutex", func(b *testing.B) {
		var c RWMutexCounter
		c.Increment()
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				sinkInt64 = c.Get()
			}
		})
	})

	b.Run("CAS_Spin", func(b *testing.B) {
		var c CASSpinCounter
		c.Increment()
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				sinkInt64 = c.Get()
			}
		})
	})


	b.Run("CAS_Backoff", func(b *testing.B) {
		var c CASBackoffCounter
		c.Increment()
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				sinkInt64 = c.Get()
			}
		})
	})

	for _, buf := range []int{0, 1, 256, 4096} {
		b.Run("Channel_buf"+strconv.Itoa(buf), func(b *testing.B) {
			c := NewChannelCounter(buf)
			b.Cleanup(c.Close)
			c.Increment()
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				reply := make(chan int64, 1)
				for pb.Next() {
					sinkInt64 = c.Get(reply)
				}
			})
		})
	}
}

func BenchmarkCounters_Mixed90_10(b *testing.B) {
	b.Run("Atomic", func(b *testing.B) {
		var c AtomicCounter
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				i++
				if i%10 == 0 {
					sinkInt64 = c.Get()
				} else {
					c.Increment()
				}
			}
		})
	})

	b.Run("Mutex", func(b *testing.B) {
		var c MutexCounter
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				i++
				if i%10 == 0 {
					sinkInt64 = c.Get()
				} else {
					c.Increment()
				}
			}
		})
	})

	b.Run("RWMutex", func(b *testing.B) {
		var c RWMutexCounter
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				i++
				if i%10 == 0 {
					sinkInt64 = c.Get()
				} else {
					c.Increment()
				}
			}
		})
	})

	b.Run("CAS_Spin", func(b *testing.B) {
		var c CASSpinCounter
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				i++
				if i%10 == 0 {
					sinkInt64 = c.Get()
				} else {
					c.Increment()
				}
			}
		})
	})

	b.Run("CAS_Backoff", func(b *testing.B) {
		var c CASBackoffCounter
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				i++
				if i%10 == 0 {
					sinkInt64 = c.Get()
				} else {
					c.Increment()
				}
			}
		})
	})

	for _, buf := range []int{0, 1, 256, 4096} {
		b.Run("Channel_buf"+strconv.Itoa(buf), func(b *testing.B) {
			c := NewChannelCounter(buf)
			b.Cleanup(c.Close)
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				reply := make(chan int64, 1)
				i := 0
				for pb.Next() {
					i++
					if i%10 == 0 {
						sinkInt64 = c.Get(reply)
					} else {
						c.Increment()
					}
				}
			})
		})
	}
}

func BenchmarkSPSC(b *testing.B) {
	// Baseline: buffered channel used as SPSC queue (non-blocking send/recv
	// with Gosched, same backpressure strategy as the ring buffer)
	for _, size := range []int{256, 4096, 65536} {
		b.Run("Chan_buf"+strconv.Itoa(size), func(b *testing.B) {
			ch := make(chan int64, size)
			b.ReportAllocs()
			b.ResetTimer()

			done := make(chan struct{})
			go func() {
				for i := 0; i < b.N; i++ {
					for {
						select {
						case ch <- int64(i):
							goto sent
						default:
							runtime.Gosched()
						}
					}
				sent:
				}
				close(done)
			}()

			for i := 0; i < b.N; i++ {
				for {
					select {
					case <-ch:
						goto recvd
					default:
						runtime.Gosched()
					}
				}
			recvd:
			}
			<-done
		})
	}

	// SPSC Ring Buffer
	for _, size := range []int{256, 4096, 65536} {
		b.Run("RingBuffer_size"+strconv.Itoa(size), func(b *testing.B) {
			rb := NewSPSCRingBuffer(size)
			b.ReportAllocs()
			b.ResetTimer()

			done := make(chan struct{})
			go func() {
				for i := 0; i < b.N; i++ {
					for !rb.Push(int64(i)) {
						runtime.Gosched()
					}
				}
				close(done)
			}()

			for i := 0; i < b.N; i++ {
				for {
					if _, ok := rb.Pop(); ok {
						break
					}
					runtime.Gosched()
				}
			}
			<-done
		})
	}
}

var mapKeys = func() []string {
	keys := make([]string, 1024)
	for i := range keys {
		keys[i] = "key" + strconv.Itoa(i)
	}
	return keys
}()

func BenchmarkMap_StoreOnly(b *testing.B) {
	b.Run("SyncMap", func(b *testing.B) {
		var s SyncMapStore
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				s.Store(mapKeys[i&1023], int64(i))
				i++
			}
		})
	})

	b.Run("MutexMap", func(b *testing.B) {
		s := NewMutexMap()
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				s.Store(mapKeys[i&1023], int64(i))
				i++
			}
		})
	})
}

func BenchmarkMap_LoadOnly(b *testing.B) {
	// Pre-populate
	var sm SyncMapStore
	mm := NewMutexMap()
	for i := 0; i < 1024; i++ {
		sm.Store(mapKeys[i], int64(i))
		mm.Store(mapKeys[i], int64(i))
	}

	b.Run("SyncMap", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				sinkInt64, _ = sm.Load(mapKeys[i&1023])
				i++
			}
		})
	})

	b.Run("MutexMap", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				sinkInt64, _ = mm.Load(mapKeys[i&1023])
				i++
			}
		})
	})
}

func BenchmarkSemaphore(b *testing.B) {
	for _, maxConc := range []int{1, 4, 16, 64} {
		b.Run("max"+strconv.Itoa(maxConc), func(b *testing.B) {
			sem := NewSemaphore(maxConc)
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					sem.Acquire()
					sem.Release()
				}
			})
		})
	}
}

func BenchmarkSpawnAndWait(b *testing.B) {
	for _, workers := range []int{1, 4, 16, 64} {
		b.Run("workers"+strconv.Itoa(workers), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				WaitGroupWork(workers, func() {})
			}
		})
	}
}

func BenchmarkPool(b *testing.B) {
	b.Run("Serial", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			PoolGetPut()
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				PoolGetPut()
			}
		})
	})
}

func BenchmarkOnce(b *testing.B) {
	b.Run("Serial", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var o OnceInit
			o.Init()
			sinkInt64 = o.Get()
		}
	})

	b.Run("Contended", func(b *testing.B) {
		var o OnceInit
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				o.Init()
				sinkInt64 = o.Get()
			}
		})
	})

	b.Run("FastPath_AfterInit", func(b *testing.B) {
		var o OnceInit
		o.Init() // init before benchmark starts
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				o.Init() // fast path: already initialized
				sinkInt64 = o.Get()
			}
		})
	})
}

// ===========================================================================
// Contention-Reduction Algorithm Benchmarks
// ===========================================================================

// ---------------------------------------------------------------------------
// BenchmarkReduction_IncOnly: write-heavy contention comparison
// ---------------------------------------------------------------------------

func BenchmarkReduction_IncOnly(b *testing.B) {
	b.Run("Sharded_8", func(b *testing.B) {
		c := NewShardedCounter(8)
		var workerID atomic.Int64
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			shard := int(workerID.Add(1))
			for pb.Next() {
				c.Increment(shard)
			}
		})
	})

	b.Run("Sharded_NumCPU", func(b *testing.B) {
		c := NewShardedCounter(runtime.NumCPU())
		var workerID atomic.Int64
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			shard := int(workerID.Add(1))
			for pb.Next() {
				c.Increment(shard)
			}
		})
	})

	b.Run("PerCPU", func(b *testing.B) {
		c := NewPerCPUCounter()
		var workerID atomic.Int64
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			slot := int(workerID.Add(1))
			for pb.Next() {
				c.Increment(slot)
			}
		})
	})

	b.Run("FlatCombining", func(b *testing.B) {
		var c FlatCombiningCounter
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				c.Increment()
			}
		})
	})

	for _, thresh := range []int64{16, 64, 256} {
		b.Run("LocalBuffered_t"+strconv.FormatInt(thresh, 10), func(b *testing.B) {
			c := NewLocalBufferedCounter(thresh)
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				var local int64
				for pb.Next() {
					c.IncrementBuffered(&local)
				}
				c.Flush(&local)
			})
		})
	}

	for _, buf := range []int{256, 4096} {
		b.Run("MPSC_buf"+strconv.Itoa(buf), func(b *testing.B) {
			a := NewMPSCAggregator(buf)
			b.Cleanup(a.Close)
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					a.Increment()
				}
			})
		})
	}

	b.Run("RCU", func(b *testing.B) {
		c := NewRCUCounter()
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				c.Increment()
			}
		})
	})

	for _, rate := range []int64{10, 100, 1000} {
		b.Run("Approx_rate"+strconv.FormatInt(rate, 10), func(b *testing.B) {
			c := NewApproxCounter(rate)
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				seed := uint64(12345)
				for pb.Next() {
					c.Increment(&seed)
				}
			})
		})
	}
}

// ---------------------------------------------------------------------------
// BenchmarkReduction_GetOnly: read contention for reduction algorithms
// ---------------------------------------------------------------------------

func BenchmarkReduction_GetOnly(b *testing.B) {
	b.Run("Sharded_8", func(b *testing.B) {
		c := NewShardedCounter(8)
		c.Increment(0) // seed
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				sinkInt64 = c.Get()
			}
		})
	})

	b.Run("PerCPU", func(b *testing.B) {
		c := NewPerCPUCounter()
		c.Increment(0)
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				sinkInt64 = c.Get()
			}
		})
	})

	b.Run("RCU", func(b *testing.B) {
		c := NewRCUCounter()
		c.Increment()
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				sinkInt64 = c.Get()
			}
		})
	})

	b.Run("Approx_rate100", func(b *testing.B) {
		c := NewApproxCounter(100)
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				sinkInt64 = c.Get()
			}
		})
	})
}

// ---------------------------------------------------------------------------
// BenchmarkDisruptorMPSC: multi-producer single-consumer throughput
// ---------------------------------------------------------------------------

func BenchmarkDisruptorMPSC(b *testing.B) {
	for _, size := range []int{256, 4096, 65536} {
		b.Run("size"+strconv.Itoa(size), func(b *testing.B) {
			d := NewDisruptorMPSC(size)
			b.ReportAllocs()
			b.ResetTimer()

			// 4 producers
			producers := 4
			perProducer := b.N / producers
			var wg sync.WaitGroup
			wg.Add(producers)
			for p := 0; p < producers; p++ {
				go func() {
					defer wg.Done()
					for i := 0; i < perProducer; i++ {
						for !d.Publish(int64(i)) {
							runtime.Gosched()
						}
					}
				}()
			}

			total := producers * perProducer
			for i := 0; i < total; i++ {
				for {
					if _, ok := d.Consume(); ok {
						break
					}
					runtime.Gosched()
				}
			}
			wg.Wait()
		})
	}
}